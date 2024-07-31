package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	common "test.com/project-common"
	"test.com/project-user/pkg/dao"
	"test.com/project-user/pkg/repo"
	"time"
)

type HandlerUser struct { // 可以理解为依赖，就是handleruser 依赖实现catche接口的结构体（对象）
	cache repo.Cache
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{
		cache: dao.Rc}
}

func (h *HandlerUser) getCaptcha(ctx *gin.Context) {
	rsp := &common.Result{}
	// 1. 获取参数
	mobile := ctx.PostForm("mobile")
	// 2. 校验参数
	if !common.VerifyMobile(mobile) {
		ctx.JSON(http.StatusOK, rsp.Fail(000000, "手机号不合法")) // model.NoLegalMobile
		return
	}
	// 3. 生成验证码（随机4位1000-9999 或者 6位 100000-999999） 因为在线的验证码服务不给个人开放，这样替代一下；
	code := "123456"
	// 4. 调用短信平台（第三方，放入go协程中，接口可以快速响应）
	go func() {
		time.Sleep(2 * time.Second)
		zap.L().Info("短信平台调用成功，发送短信")

		// redis 假设后续缓存可能存到mysql中，也可能存到mongoDB中，也可能存在memcache中，
		// 如果存的方式变了，就要修改这里代码了，引入（repo）类Service层 和 （dao）DAO层 的概念。
		// 5. 存储验证码 redis 当中，过期时间为15分钟
		// redis.Set("REGISTER_"+mobile, code)

		c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := h.cache.Put(c, "REGISTER_"+mobile, code, 15*time.Minute)
		if err != nil {
			log.Println("验证码存入redis出错：", err)
		}
		log.Printf("将手机号和验证码存入redis成功：REGISTER_%s : %s", mobile, code)
	}()
	ctx.JSON(200, rsp.Success(code))
}
