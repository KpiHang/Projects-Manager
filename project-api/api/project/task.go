package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"test.com/project-api/pkg/model"
	"test.com/project-api/pkg/model/pro"
	"test.com/project-api/pkg/model/tasks"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/task"
	"time"
)

type HandlerTask struct {
}

func NewTask() *HandlerTask {
	return &HandlerTask{}
}

func (t HandlerTask) taskStages(c *gin.Context) {
	// 1. 获取参数，可能校验参数
	// 2. 调用对应的grpc服务
	// 3. 返回响应。

	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 1. 获取参数，可能校验参数
	projectCode := c.PostForm("projectCode")
	page := &model.Page{}
	page.Bind(c)
	// 2. 调用对应的grpc服务
	msg := &task.TaskReqMessage{
		MemberId:    c.GetInt64("memberId"),
		ProjectCode: projectCode,
		Page:        page.Page,
		PageSize:    page.PageSize,
	}
	stages, err := TaskServiceClient.TaskStages(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	// 3. 返回响应。
	var list []*tasks.TaskStagesResp
	copier.Copy(&list, stages.List)
	if list == nil {
		list = []*tasks.TaskStagesResp{}
	}
	for _, v := range list {
		v.TasksLoading = true  //任务加载状态
		v.FixedCreator = false //添加任务按钮定位
		v.ShowTaskCard = false //是否显示创建卡片
		v.Tasks = []int{}
		v.DoneTasks = []int{}
		v.UnDoneTasks = []int{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  list,
		"total": stages.Total,
		"page":  page.Page,
	}))
}

func (t HandlerTask) memberProjectList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 1. 获取参数，可能校验参数
	projectCode := c.PostForm("projectCode")
	page := &model.Page{}
	page.Bind(c)
	// 2. 调用对应的grpc服务
	msg := &task.TaskReqMessage{
		MemberId:    c.GetInt64("memberId"),
		ProjectCode: projectCode,
		Page:        page.Page,
		PageSize:    page.PageSize,
	}
	resp, err := TaskServiceClient.MemberProjectList(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	// 3. 返回响应。
	var list []*pro.MemberProjectResp
	copier.Copy(&list, resp.List)
	if list == nil {
		list = []*pro.MemberProjectResp{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  list,
		"total": resp.Total,
		"page":  page.Page,
	}))
}

func (t HandlerTask) taskList(c *gin.Context) {
	result := &common.Result{}
	stageCode := c.PostForm("stageCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	list, err := TaskServiceClient.TaskList(ctx, &task.TaskReqMessage{StageCode: stageCode})
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var taskDisplayList []*tasks.TaskDisplay
	copier.Copy(&taskDisplayList, list.List)
	if taskDisplayList == nil {
		taskDisplayList = []*tasks.TaskDisplay{}
	}
	// 返回给前端的数据，一定不能是none nil
	for _, v := range taskDisplayList {
		if v.Tags == nil {
			v.Tags = []int{}
		}
		if v.ChildCount == nil {
			v.ChildCount = []int{}
		}
	}
	c.JSON(http.StatusOK, result.Success(taskDisplayList))
}
