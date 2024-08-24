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
	group.POST("/project/getLogBySelfProject", h.getLogBySelfProject)

	t := NewTask()
	group.POST("/task_stages", t.taskStages)
	group.POST("/project_member/index", t.memberProjectList)
	group.POST("/task_stages/tasks", t.taskList)
	group.POST("/task/save", t.saveTask)
	group.POST("/task/sort", t.taskSort)
	group.POST("/task/selfList", t.myTaskList)
	group.POST("/task/read", t.readTask)
	group.POST("/task_member", t.listTaskMember)
	group.POST("/task/taskLog", t.taskLog)
	group.POST("/task/_taskWorkTimeList", t.taskWorkTimeList)
	group.POST("/task/saveTaskWorkTime", t.saveTaskWorkTime)
	group.POST("/file/uploadFiles", t.uploadFiles)
	group.POST("/task/taskSources", t.taskSources)
	group.POST("/task/createComment", t.createComment)

	a := NewAccount()
	group.POST("/account", a.account)

	d := NewDepartment()
	group.POST("/department", d.department)
	group.POST("/department/save", d.save)
	group.POST("/department/read", d.read)

	auth := NewAuth()
	group.POST("/auth", auth.authList)

	menu := NewMenu()
	group.POST("/menu/menu", menu.menuList)

}
