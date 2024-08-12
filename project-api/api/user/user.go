package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"test.com/project-api/api/rpc"
	"test.com/project-api/pkg/model/user"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/user/login"
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
	captchaRsp, err := rpc.LoginServiceClient.GetCaptcha(c, &login.CaptchaMessage{Mobile: mobile})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		ctx.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	ctx.JSON(http.StatusOK, result.Success(captchaRsp.Code)) // 由api网关响应给客户端
}

func (h *HandlerUser) Register(ctx *gin.Context) {
	// 1. 接收参数；需要有一个参数的模型（结构体、Model）
	// 2. 校验参数；参数是否合法；
	// 3. 调用user 注册的grpc服务；
	// 4. 返回响应

	gatewayResponse := &common.Result{}
	// 1. 接收参数；需要有一个参数的模型（结构体、Model）
	var req user.RegisterReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, gatewayResponse.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	// 2. 校验参数；参数是否合法；
	if err := req.Verify(); err != nil {
		ctx.JSON(http.StatusOK, gatewayResponse.Fail(http.StatusBadRequest, err.Error()))
		return
	}
	// 3. 调用user 注册的grpc服务；
	c, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()
	msg := &login.RegisterMessage{}
	err = copier.Copy(msg, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gatewayResponse.Fail(http.StatusBadRequest, "copy 参数格式有误"))
		return
	}
	_, err = rpc.LoginServiceClient.Register(c, msg) // 在user 模块中写注册相关的grpc服务；
	if err != nil {
		code, msg := errs.ParseGrpcError(err) // grpc 服务返回的code msg
		ctx.JSON(http.StatusOK, gatewayResponse.Fail(code, msg))
		return
	}
	// 4. 返回响应
	ctx.JSON(http.StatusOK, gatewayResponse.Success("")) // 由api网关响应给客户端
}

func (h *HandlerUser) Login(ctx *gin.Context) {
	// 1. 接收参数；需要有一个参数的模型（结构体、Model）
	// 2. 调用user grpc 完成登陆
	// 3. 返回响应

	gatewayResponse := &common.Result{}
	// 1. 接收参数；需要有一个参数的模型（结构体、Model）
	var req user.LoginReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, gatewayResponse.Fail(http.StatusBadRequest, "参数格式有误"))
		return
	}
	// 2. 调用user grpc 完成登陆
	c, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()
	msg := &login.LoginMessage{}
	err = copier.Copy(msg, req)
	if err != nil {
		ctx.JSON(http.StatusOK, gatewayResponse.Fail(http.StatusBadRequest, "copy 参数格式有误"))
		return
	}
	loginRsp, err := rpc.LoginServiceClient.Login(c, msg) // 在user 模块中写注册相关的grpc服务；
	if err != nil {
		code, msg := errs.ParseGrpcError(err) // grpc 服务返回的code msg
		ctx.JSON(http.StatusOK, gatewayResponse.Fail(code, msg))
		return
	}
	rsp := &user.LoginRsp{}
	err = copier.Copy(rsp, loginRsp)
	// 3. 返回响应
	ctx.JSON(http.StatusOK, gatewayResponse.Success(rsp)) // 由api网关响应给客户端
}

func (h *HandlerUser) myOrgList(c *gin.Context) {
	result := &common.Result{}
	memberIdStr, _ := c.Get("memberId")
	memberId := memberIdStr.(int64)
	list, err2 := rpc.LoginServiceClient.MyOrgList(context.Background(), &login.UserMessage{MemId: memberId})

	if err2 != nil {
		code, msg := errs.ParseGrpcError(err2)
		c.JSON(http.StatusOK, result.Fail(code, msg))
		return
	}
	if list.OrganizationList == nil {
		c.JSON(http.StatusOK, result.Success([]*user.OrganizationList{}))
		return
	}
	var orgs []*user.OrganizationList
	copier.Copy(&orgs, list.OrganizationList)
	c.JSON(http.StatusOK, result.Success(orgs))
}
