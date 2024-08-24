package midd

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"test.com/project-api/api/rpc"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/user/login"
	"time"
)

// 获取ip函数
func GetIp(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}

func TokenVerify() func(c *gin.Context) {
	return func(c *gin.Context) {
		result := &common.Result{}
		// 1. 从header中获取token
		// 2. 调用user service进行token认证；
		// 3. 如果认证通过，将信息放入gin的上下文；失败就返回未登录；

		// 1. 从header中获取token
		token := c.GetHeader("Authorization")
		// 2. 调用user service进行token认证；
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		ip := GetIp(c)
		// todo: 先查询node表，确认不使用登陆控制的接口，就不做登陆认证了。
		verifyRsp, err := rpc.LoginServiceClient.TokenVerify(ctx, &login.LoginMessage{Token: token, Ip: ip})
		if err != nil {
			code, msg := errs.ParseGrpcError(err)
			c.JSON(http.StatusOK, result.Fail(code, msg))
			c.Abort() // 中止当前的请求处理流程; Gin 将停止执行后续的中间件或处理器函数，并立即返回响应给客户端。
			return
		}
		// 3. 如果认证通过，将信息（按需）放入gin的上下文；失败就返回未登录；
		c.Set("memberId", verifyRsp.Member.Id) // 用于在请求的上下文中存储一个键值对
		c.Set("memberName", verifyRsp.Member.Name)
		c.Set("organizationCode", verifyRsp.Member.OrganizationCode)
		c.Next() // 在中间件中明确调用下一个处理器
	}
}
