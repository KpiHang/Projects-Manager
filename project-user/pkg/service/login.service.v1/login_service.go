package login_service_v1

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"log"
	"strconv"
	"strings"
	common "test.com/project-common"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/jwts"
	"test.com/project-common/tms"
	"test.com/project-grpc/user/login"
	"test.com/project-user/config"
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
	login.UnimplementedLoginServiceServer                       // grpc 里的
	cache                                 repo.Cache            // repo里定义接口
	memberRepo                            repo.MemberRepo       // repo里定义接口
	organizationRepo                      repo.OrganizationRepo // repo里定义接口
	transaction                           tran.Transaction      // 事务操作接口
}

// NewLoginService 因为catche字段是接口，构造一下把链接redis后的cache放进来；
func NewLoginService() *LoginService {
	return &LoginService{
		cache:            dao.Rc,                   // dao中实现接口，按照repo接口规则和数据库交互
		memberRepo:       dao.NewMemberDao(),       // dao中实现接口，按照repo接口规则和数据库交互
		organizationRepo: dao.NewOrganizationDao(), // dao中实现接口，按照repo接口规则和数据库交互
		transaction:      dao.NewTransactionImpl(), // dao中实现接口，按照repo接口规则和数据库交互
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

func (ls *LoginService) Login(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	c := context.Background()
	// 1. 去数据库查询，账号密码是否正确；
	pwd := encrypts.Md5(msg.Password)
	mem, err := ls.memberRepo.FindMember(c, msg.Account, pwd)
	if err != nil {
		zap.L().Error("Login DB error, ", zap.Error(err)) // 非业务错误；
		return nil, errs.GrpcError(model.DBError)
	}
	// todo: 如果查询为空，mem 并不是是nil，而是一个零值 &member.Member{}
	if mem == nil {
		return nil, errs.GrpcError(model.AccountOrPwdError)
	}
	memMessage := &login.MemberMessage{} // grpc服务的响应实体（之一）
	err = copier.Copy(memMessage, mem)
	memMessage.Code, _ = encrypts.EncryptInt64(mem.Id, model.AESKey) // 加密id
	memMessage.LastLoginTime = tms.FormatByMill(mem.LastLoginTime)
	memMessage.CreateTime = tms.FormatByMill(mem.CreateTime)
	// 2. 根据用户id 查组织；
	orgs, err := ls.organizationRepo.FindOrganizationByMemberId(c, mem.Id)
	if err != nil {
		zap.L().Error("Login DB error, ", zap.Error(err)) // 非业务错误，
		return nil, errs.GrpcError(model.DBError)
	}

	var orgsMessage []*login.OrganizationMessage // grpc服务的响应实体（之一）
	err = copier.Copy(&orgsMessage, orgs)
	for _, org := range orgsMessage {
		org.Code, _ = encrypts.EncryptInt64(org.Id, model.AESKey) // 加密组织的id
		org.OwnerCode = memMessage.Code
		org.CreateTime = tms.FormatByMill(organization.ToMap(orgs)[org.Id].CreateTime)
	}

	// 3. 用jwt生成token
	memIdStr := strconv.FormatInt(mem.Id, 10)
	exp := time.Duration(config.Conf.JwtConfig.AccessExp*3600*24) * time.Second
	rexp := time.Duration(config.Conf.JwtConfig.RefreshExp*3600*24) * time.Second

	token := jwts.CreateToken(
		memIdStr,
		exp,
		config.Conf.JwtConfig.AccessSecret,
		rexp,
		config.Conf.JwtConfig.RefreshSecret,
	)

	tokenList := &login.TokenMessage{
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		AccessTokenExp: token.AccessExp,
		TokenType:      "bearer",
	}

	// 4. 返回响应
	return &login.LoginResponse{
		Member:           memMessage,
		OrganizationList: orgsMessage,
		TokenList:        tokenList,
	}, nil
}

func (ls *LoginService) TokenVerify(ctx context.Context, msg *login.LoginMessage) (*login.LoginResponse, error) {
	token := msg.Token
	if strings.Contains(token, "bearer") {
		token = strings.ReplaceAll(token, "bearer ", "") // 前端传过来的token 前面带了 bearer ;
	}
	parseToken, err := jwts.ParseToken(token, config.Conf.JwtConfig.AccessSecret)
	if err != nil {
		zap.L().Error("TokenVerify ParseToken err", zap.Error(err))
		return nil, errs.GrpcError(model.NoLogin)
	}

	// 数据库查询，优化点（todo），登陆之后应该把用户信息缓存起来；
	id, _ := strconv.ParseInt(parseToken, 10, 64)
	memberById, err := ls.memberRepo.FindMemberById(context.Background(), id)
	if err != nil {
		zap.L().Error("TokenVerify FindMemberById DB error, ", zap.Error(err)) // 非业务错误，
		return nil, errs.GrpcError(model.DBError)
	}
	memMessage := &login.MemberMessage{} // grpc服务的响应实体（之一）
	copier.Copy(memMessage, memberById)
	memMessage.Code, _ = encrypts.EncryptInt64(memberById.Id, model.AESKey) // 加密id

	return &login.LoginResponse{Member: memMessage}, nil
}

func (l *LoginService) MyOrgList(ctx context.Context, msg *login.UserMessage) (*login.OrgListResponse, error) {
	memId := msg.MemId
	orgs, err := l.organizationRepo.FindOrganizationByMemberId(ctx, memId)
	if err != nil {
		zap.L().Error("MyOrgList FindOrganizationByMemId err", zap.Error(err))
		return nil, errs.GrpcError(model.DBError)
	}
	var orgsMessage []*login.OrganizationMessage
	err = copier.Copy(&orgsMessage, orgs)
	for _, org := range orgsMessage {
		org.Code, _ = encrypts.EncryptInt64(org.Id, model.AESKey)
	}
	return &login.OrgListResponse{OrganizationList: orgsMessage}, nil
}
