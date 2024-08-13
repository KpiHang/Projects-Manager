package project

import (
	"github.com/gin-gonic/gin"
	"log"
	"test.com/project-api/api/midd"
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
	group := r.Group("/project/index")
	group.Use(midd.TokenVerify()) // 这个组用中间件；
	group.POST("", h.index)

	group1 := r.Group("/project/project")
	group1.Use(midd.TokenVerify())
	group1.POST("/selfList", h.myProjectList) // 用id获取我的项目list
	group1.POST("", h.myProjectList)          // 用id获取 select对应类型的项目list  (表单多传一个selectBy)
}
