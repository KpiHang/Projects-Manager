package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

// AccountRepo 其实是MemberAccountRepo
type AccountRepo interface {
	FindList(ctx context.Context, condition string, organizationCode int64, departmentCode int64, page int64, pageSize int64) ([]*data.MemberAccount, int64, error)
	FindByMemberId(ctx context.Context, memberId int64) (*data.MemberAccount, error)
}
