package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"test.com/project-api/pkg/model"
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
	c.JSON(http.StatusOK, result.Success(indexResponse.Menus))
}

func (p HandlerProject) myProjectList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 1. 获取参数；
	memberIdStr, _ := c.Get("memberId") // 自定义中间件放进去的
	memberId := memberIdStr.(int64)
	page := &model.Page{}
	page.Bind(c)
	msg := &project.ProjectRpcMessage{MemberId: memberId, Page: page.Page, PageSize: page.PageSize}
	myProjectResponse, err := ProjectServiceClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	if myProjectResponse.Pm == nil {
		myProjectResponse.Pm = []*project.ProjectMessage{}
	}

	var pms []*pro.ProjectAndMember
	copier.Copy(&pms, myProjectResponse.Pm)

	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms,
		"total": myProjectResponse.Total,
	}))
}
