package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"test.com/project-api/pkg/model"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/account"
	"time"
)

type HandlerAccount struct {
}

func (a *HandlerAccount) account(c *gin.Context) {
	// 接收请求参数；（一些参数校验可以放在api里）
	// 调用project模块的grpc服务，查询账户列表；
	// 范湖响应数据；

	result := &common.Result{}
	var req *model.AccountReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// 接收请求参数；（一些参数校验可以放在api里）
	msg := &account.AccountReqMessage{
		MemberId:         c.GetInt64("memberId"),
		OrganizationCode: c.GetString("organizationCode"),
		Page:             int64(req.Page),
		PageSize:         int64(req.PageSize),
		SearchType:       int32(req.SearchType),
		DepartmentCode:   req.DepartmentCode,
	}
	response, err := AccountServiceClient.Account(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	// 返回数据
	var list []*model.MemberAccount
	copier.Copy(&list, response.AccountList)
	if list == nil {
		list = []*model.MemberAccount{}
	}
	var authList []*model.ProjectAuth
	copier.Copy(&authList, response.AuthList)
	if authList == nil {
		authList = []*model.ProjectAuth{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total":    response.Total,
		"page":     req.Page,
		"list":     list,
		"authList": authList,
	}))

}

func NewAccount() *HandlerAccount {
	return &HandlerAccount{}
}
