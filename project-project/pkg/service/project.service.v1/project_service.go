package project_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-common/tms"
	"test.com/project-grpc/project"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data/menu"
	"test.com/project-project/internal/data/pro"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
)

type ProjectService struct {
	project.UnimplementedProjectServiceServer                  // grpc 里的
	cache                                     repo.Cache       // repo里定义接口
	transaction                               tran.Transaction // 事务操作接口
	menuRepo                                  repo.MenuRepo
	projectRepo                               repo.ProjectRepo
}

func NewProjectService() *ProjectService {
	return &ProjectService{
		cache:       dao.Rc,
		transaction: dao.NewTransactionImpl(),
		menuRepo:    dao.NewMenuDao(),
		projectRepo: dao.NewProjectDao(),
	}
}

func (p *ProjectService) Index(context.Context, *project.IndexMessage) (*project.IndexResponse, error) {
	pms, err := p.menuRepo.FindMenus(context.Background())
	if err != nil {
		zap.L().Error("Index FindMenus DB error, ", zap.Error(err)) // 非业务错误；
		return nil, errs.GrpcError(model.DBError)
	}
	childs := menu.CovertChild(pms)
	var mms []*project.MenuMessage
	copier.Copy(&mms, childs)
	return &project.IndexResponse{Menus: mms}, nil
}

func (p *ProjectService) FindProjectByMemId(ctx context.Context, msg *project.ProjectRpcMessage) (*project.MyProjectResponse, error) {
	memberId := msg.MemberId
	page := msg.Page
	pageSize := msg.PageSize
	var (
		pms   []*pro.ProjectAndMember
		total int64
		err   error
	)

	if msg.SelectBy == "" || msg.SelectBy == "my" {
		// 要从数据库中query了，就调一个repo；
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, "", page, pageSize)
	}

	if msg.SelectBy == "archive" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, "and archive = 1", page, pageSize)
	}

	if msg.SelectBy == "deleted" {
		pms, total, err = p.projectRepo.FindProjectByMemId(ctx, memberId, "and deleted = 1", page, pageSize)
	}

	if msg.SelectBy == "collect" { // 用到用户项目收藏表了，用新的方法。
		pms, total, err = p.projectRepo.FindCollectProjectByMemId(ctx, memberId, page, pageSize)
	}

	if err != nil {
		zap.L().Error("project FindProjectByMemId DB error, ", zap.Error(err)) // 非业务错误；
		return nil, errs.GrpcError(model.DBError)
	}
	if pms == nil { // 如果没值的话，返回一个默认的数据；
		return &project.MyProjectResponse{Pm: []*project.ProjectMessage{}, Total: total}, nil
	}

	var pmm []*project.ProjectMessage
	copier.Copy(&pmm, pms)
	for _, v := range pmm {
		v.Code, _ = encrypts.EncryptInt64(v.Id, model.AESKey)
		pam := pro.ToMap(pms)[v.Id]
		v.AccessControlType = pam.GetAccessControlType()
		v.OrganizationCode, _ = encrypts.EncryptInt64(pam.OrganizationCode, model.AESKey)
		v.JoinTime = tms.FormatByMill(pam.JoinTime)
		v.OwnerName = msg.MemberName
		v.Order = int32(pam.Sort)
		v.CreateTime = tms.FormatByMill(pam.CreateTime)
	}
	return &project.MyProjectResponse{Pm: pmm, Total: total}, nil
}
