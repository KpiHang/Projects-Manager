package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	common "test.com/project-common"
	"test.com/project-common/errs"
	LoginServiceV1 "test.com/project-user/pkg/service/login.service.v1"
	"time"
)

type HandlerUser struct { // 可以理解为依赖，就是handleruser 依赖实现catche接口的结构体（对象）

}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (h *HandlerUser) getCaptcha(ctx *gin.Context) {
	result := &common.Result{}
	mobile := ctx.PostForm("mobile") // 客户端给API发请求，携带mobile参数；
	c, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle() // ↓ rpc 调用user service server
	captchaRsp, err := LoginServiceClient.GetCaptcha(c, &LoginServiceV1.CaptchaMessage{Mobile: mobile})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(captchaRsp.Code)) // 由api网关响应给客户端
}
