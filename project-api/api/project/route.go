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
	group := r.Group("/project")
	group.Use(midd.TokenVerify()) // 这个组用中间件；
	group.POST("/index", h.index)
	group.POST("/project/selfList", h.myProjectList) // 用id获取我的项目list
	group.POST("/project", h.myProjectList)          // 用id获取 select对应类型的项目list  (表单多传一个selectBy)
	group.POST("/project_template", h.projectTemplate)
	group.POST("/project/save", h.projectSave)
	group.POST("/project/read", h.readProject)
	group.POST("/project/recycle", h.recycleProject)         // 移入回收站
	group.POST("/project/recovery", h.recoveryProject)       // 移出回收站
	group.POST("/project_collect/collect", h.collectProject) // 收藏项目
	group.POST("/project/edit", h.editProject)               // 编辑项目
}
