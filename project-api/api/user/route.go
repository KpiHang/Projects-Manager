package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/api/midd"
	"test.com/project-api/api/rpc"
	"test.com/project-api/router"
)

func init() {
	log.Println("init user router")
	router.AddRouter(&RouterUser{})
}

type RouterUser struct {
}

func (*RouterUser) Register(r *gin.Engine) {
	// 初始化grpc客户端的连接，链接user service server，在rpc.go中完成；
	rpc.InitRpcUserClient()
	h := NewHandlerUser()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/login/register", h.Register)
	r.POST("/project/login", h.Login)

	org := r.Group("/project/organization")
	org.Use(midd.TokenVerify())
	org.POST("/_getOrgList", h.myOrgList)
}
