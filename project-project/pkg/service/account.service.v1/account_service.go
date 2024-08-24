package account_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"test.com/project-common/encrypts"
	"test.com/project-common/errs"
	"test.com/project-grpc/account"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/database/tran"
	"test.com/project-project/internal/domain"
	"test.com/project-project/internal/repo"
)

type AccountService struct {
	account.UnimplementedAccountServiceServer                       // grpc 里的
	cache                                     repo.Cache            // repo里定义接口
	transaction                               tran.Transaction      // 事务操作接口
	accountDomain                             *domain.AccountDomain // 第一个直接写领域层的代码，而不是改造
	projectAuthDomain                         *domain.ProjectAuthDomain
}

func New() *AccountService {
	return &AccountService{
		cache:             dao.Rc,
		transaction:       dao.NewTransactionImpl(),
		accountDomain:     domain.NewAccountDomain(),
		projectAuthDomain: domain.NewProjectAuthDomain(),
	}
}

func (a *AccountService) Account(ctx context.Context, msg *account.AccountReqMessage) (*account.AccountResponse, error) {
	// 1. 去 account 表里查询 account
	// 2. 去 project_auth 表里查询 authList

	// 1. 去 account 表里查询 account
	accountList, total, err := a.accountDomain.AccountList( // 调用domain,交给account domain去作，专人专事。
		msg.OrganizationCode,
		msg.MemberId,
		msg.Page,
		msg.PageSize,
		msg.DepartmentCode,
		msg.SearchType)
	if err != nil {
		return nil, errs.GrpcError(err)
	}
	// 2. 去 project_auth 表里查询 authList
	authList, err := a.projectAuthDomain.AuthList(encrypts.DecryptNoErr(msg.OrganizationCode))
	if err != nil {
		return nil, errs.GrpcError(err)
	}

	var maList []*account.MemberAccount
	copier.Copy(&maList, accountList)
	var prList []*account.ProjectAuth
	copier.Copy(&prList, authList)
	return &account.AccountResponse{
		AccountList: maList,
		AuthList:    prList,
		Total:       total,
	}, nil
}
