package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"test.com/project-api/pkg/model"
	"test.com/project-api/pkg/model/menu"
	"test.com/project-api/pkg/model/pro"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/project"
	"time"
)

type HandlerProject struct { // 可以理解为依赖，就是handleruser 依赖实现catche接口的结构体（对象）

}

func NewHandlerProject() *HandlerProject {
	return &HandlerProject{}
}

func (p HandlerProject) index(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &project.IndexMessage{}
	indexResponse, err := ProjectServiceClient.Index(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	menus := indexResponse.Menus
	var ms []*menu.Menu
	copier.Copy(&ms, menus)
	c.JSON(http.StatusOK, result.Success(ms))
}

func (p HandlerProject) myProjectList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 1. 获取参数；
	memberId := c.GetInt64("memberId") // 自定义中间件放进去的
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c) // 绑定页号和页size
	selectBy := c.PostForm("selectBy")
	msg := &project.ProjectRpcMessage{
		MemberId:   memberId,
		MemberName: memberName,
		SelectBy:   selectBy,
		Page:       page.Page,
		PageSize:   page.PageSize}
	myProjectResponse, err := ProjectServiceClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var pms []*pro.ProjectAndMember
	copier.Copy(&pms, myProjectResponse.Pm)
	if pms == nil {
		pms = []*pro.ProjectAndMember{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms, // 返回切片的时候，如果为空一定不能返回 nil，（导致前端报错），而是返回 []
		"total": myProjectResponse.Total,
	}))
}

func (p HandlerProject) projectTemplate(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 1. 获取参数；
	memberId := c.GetInt64("memberId") // 自定义中间件放进去的
	memberName := c.GetString("memberName")
	page := &model.Page{}
	page.Bind(c) // 绑定页号和页size
	viewTypeStr := c.PostForm("viewType")
	viewType, _ := strconv.ParseInt(viewTypeStr, 10, 64)
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		MemberName:       memberName,
		ViewType:         int32(viewType),
		Page:             page.Page,
		PageSize:         page.PageSize,
		OrganizationCode: c.GetString("organizationCode")}
	templateResponse, err := ProjectServiceClient.FindProjectTemplate(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var pms []*pro.ProjectTemplate
	copier.Copy(&pms, templateResponse.Ptm)
	if pms == nil {
		pms = []*pro.ProjectTemplate{}
	}

	for _, v := range pms {
		if v.TaskStages == nil { // 因为TaskStages也是列表，所有返回不能为nil 而是 []
			v.TaskStages = []*pro.TaskStagesOnlyName{}
		}
	}

	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms, // 返回切片的时候，如果为空一定不能返回 nil，（导致前端报错），而是返回 []
		"total": templateResponse.Total,
	}))
}

func (p HandlerProject) projectSave(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 1. 获取参数；
	memberId := c.GetInt64("memberId") // 自定义中间件放进去的
	organizationCode := c.GetString("organizationCode")

	var req *pro.SaveProjectRequest
	c.ShouldBind(&req)
	msg := &project.ProjectRpcMessage{
		MemberId:         memberId,
		OrganizationCode: organizationCode,
		TemplateCode:     req.TemplateCode,
		Name:             req.Name,
		Id:               int64(req.Id),
		Description:      req.Description,
	}
	saveProject, err := ProjectServiceClient.SaveProject(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var rsp *pro.SaveProject
	copier.Copy(&rsp, saveProject)
	c.JSON(http.StatusOK, result.Success(rsp))
}

func (p HandlerProject) readProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	memberId := c.GetInt64("memberId") // 传这个有点蠢hhh；
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	detail, err := ProjectServiceClient.FindProjectDetail(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		MemberId:    memberId,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	pd := &pro.ProjectDetail{}
	copier.Copy(pd, detail)
	c.JSON(http.StatusOK, result.Success(pd))
}

func (p HandlerProject) recycleProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectServiceClient.UpdateDeteledProject(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		Deleted:     true,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{})) // 返回值为空切片就行；
}

// recoveryProject 恢复项目，和删除项目用一个接口，传递参数  Deleted:     false,;  也可以用俩接口；
func (p HandlerProject) recoveryProject(c *gin.Context) {
	result := &common.Result{}
	projectCode := c.PostForm("projectCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := ProjectServiceClient.UpdateDeteledProject(ctx, &project.ProjectRpcMessage{
		ProjectCode: projectCode,
		Deleted:     false,
	})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	c.JSON(http.StatusOK, result.Success([]int{})) // 返回值为空切片就行；
}
