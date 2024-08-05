package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/router"
)

func init() {
	log.Println("init project router")
	router.AddRouter(&RouterProject{})
}

type RouterProject struct {
}

func (*RouterProject) Register(r *gin.Engine) {
	// 初始化grpc客户端的连接，链接user service server，在rpc.go中完成；
	InitRpcProjectClient()
	h := NewHandlerProject()
	r.POST("/project/index", h.index)
}
