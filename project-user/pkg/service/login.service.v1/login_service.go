package login_service_v1

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"log"
	common "test.com/project-common"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-grpc/user/login"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/data/organization"
	"test.com/project-user/internal/database"
	"test.com/project-user/internal/database/tran"
	"test.com/project-user/internal/repo"
	"test.com/project-user/pkg/model"
	"time"
)

type LoginService struct {
	login.UnimplementedLoginServiceServer
	cache            repo.Cache
	memberRepo       repo.MemberRepo
	organizationRepo repo.OrganizationRepo
	transaction      tran.Transaction // 事务操作接口
}

// NewLoginService 因为catche字段是接口，构造一下把链接redis后的cache放进来；
func NewLoginService() *LoginService {
	return &LoginService{
		cache:            dao.Rc,
		memberRepo:       dao.NewMemberDao(),
		organizationRepo: dao.NewOrganizationDao(),
		transaction:      dao.NewTransactionImpl(),
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *login.CaptchaMessage) (*login.CaptchaResponse, error) {
	// 1. 获取参数
	mobile := msg.Mobile
	// 2. 校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errs.GrpcError(model.NoLegalMobile)
	}
	// 3. 生成验证码（随机4位1000-9999 或者 6位 100000-999999） 因为在线的验证码服务不给个人开放，这样替代一下；
	code := "123458"
	// 4. 调用短信平台（第三方，放入go协程中，接口可以快速响应）
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("短信平台调用成功，发送短信")

		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := ls.cache.Put(c, model.RegisterRedisKey+mobile, code, 15*time.Minute)
		if err != nil {
			log.Println("验证码存入redis出错：", err)
		}
		log.Printf("将手机号和验证码存入redis成功：REGISTER_%s : %s", mobile, code)
	}()
	return &login.CaptchaResponse{Code: code}, nil
}

func (ls *LoginService) Register(ctx context.Context, msg *login.RegisterMessage) (*login.RegisterResponse, error) {
	// 1. 可以校验参数；（网关上已经初步校验过了， 可以省略）
	// 2. 校验验证码；
	// 3. 校验业务逻辑；（邮箱是否被注册，账号是否被注册，手机号是否被注册）
	// 4. 执行业务逻辑，将数据存入member表，生成一个数据，存入组织表organization
	// 5. 返回响应

	// 2. 校验验证码；
	c := context.Background()
	redisCode, err := ls.cache.Get(c, model.RegisterRedisKey+msg.Mobile) // 从redis中获取验证码
	if errors.Is(err, redis.Nil) {
		return nil, errs.GrpcError(model.CaptchaNotExist)
	}
	if err != nil {
		zap.L().Error("从redis中获取验证码出错：", zap.Error(err))
		return nil, errs.GrpcError(model.RedisError)
	}
	if redisCode != msg.Captcha {
		return nil, errs.GrpcError(model.CaptchaError)
	}
	// 3. 校验业务逻辑；（邮箱是否被注册，账号是否被注册，手机号是否被注册）
	exist, err := ls.memberRepo.IsInMemberEmail(c, msg.Email) // 邮箱
	if err != nil {
		zap.L().Error("Register DB error, ", zap.Error(err)) // 非业务错误；
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.EmailExist) // 业务错误；要返回给客户端的；
	}

	exist, err = ls.memberRepo.IsInMemberAccount(c, msg.Name) // 账号就是用户名 name
	if err != nil {
		zap.L().Error("Register DB error, ", zap.Error(err)) // 非业务错误；
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.AccountExist) // 业务错误；要返回给客户端的；
	}

	exist, err = ls.memberRepo.IsInMemberMobile(c, msg.Mobile) // 手机号
	if err != nil {
		zap.L().Error("Register DB error, ", zap.Error(err)) // 非业务错误；
		return nil, errs.GrpcError(model.DBError)
	}
	if exist {
		return nil, errs.GrpcError(model.MobileExist) // 业务错误；要返回给客户端的；
	}
	// 4. 执行业务逻辑，将数据存入member表，生成一个数据，存入组织表organization
	pwd := encrypts.Md5(msg.Password)
	mem := &member.Member{
		Account:       msg.Name,
		Password:      pwd,
		Name:          msg.Name,
		Mobile:        msg.Mobile,
		Email:         msg.Email,
		CreateTime:    time.Now().UnixMilli(),
		LastLoginTime: time.Now().UnixMilli(),
		Status:        model.Normal,
	}
	err = ls.transaction.Action(func(conn database.DbConn) error {
		err = ls.memberRepo.SaveMember(conn, c, mem)
		if err != nil {
			zap.L().Error("Register DB  save member error, ", zap.Error(err)) // 非业务错误，
			return errs.GrpcError(model.DBError)
		}
		// 用户创建完之后，还要存入组织；
		org := &organization.Organization{
			Name:       mem.Name + "个人组织",
			MemberId:   mem.Id,
			CreateTime: time.Now().UnixMilli(),
			Personal:   model.Personal,
			Avatar:     "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fc-ssl.dtstatic.com%2Fuploads%2Fblog%2F202103%2F31%2F20210331160001_9a852.thumb.1000_0.jpg&refer=http%3A%2F%2Fc-ssl.dtstatic.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=auto?sec=1673017724&t=ced22fc74624e6940fd6a89a21d30cc5",
		}
		err = ls.organizationRepo.SaveOrganization(conn, c, org)
		if err != nil {
			zap.L().Error("register SaveOrganization db err", zap.Error(err))
			return errs.GrpcError(model.DBError)
		}
		return nil
	})

	// 5. 返回响应
	return &login.RegisterResponse{}, err
}
