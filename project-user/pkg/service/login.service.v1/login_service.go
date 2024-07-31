package login_service_v1

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"log"
	common "test.com/project-common"
	"test.com/project-user/pkg/dao"
	"test.com/project-user/pkg/repo"
	"time"
)

type LoginService struct {
	UnimplementedLoginServiceServer
	cache repo.Cache
}

// NewLoginService 因为catche字段是接口，构造一下把链接redis后的cache放进来；
func NewLoginService() *LoginService {
	return &LoginService{
		cache: dao.Rc,
	}
}

func (ls *LoginService) GetCaptcha(ctx context.Context, msg *CaptchaMessage) (*CaptchaResponse, error) {
	// 1. 获取参数
	mobile := msg.Mobile
	// 2. 校验参数
	if !common.VerifyMobile(mobile) {
		return nil, errors.New("手机号不合法")
	}
	// 3. 生成验证码（随机4位1000-9999 或者 6位 100000-999999） 因为在线的验证码服务不给个人开放，这样替代一下；
	code := "123456"
	// 4. 调用短信平台（第三方，放入go协程中，接口可以快速响应）
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("短信平台调用成功，发送短信")

		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := ls.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if err != nil {
			log.Println("验证码存入redis出错：", err)
		}
		log.Printf("将手机号和验证码存入redis成功：REGISTER_%s : %s", mobile, code)
	}()
	return &CaptchaResponse{Code: code}, nil
}
