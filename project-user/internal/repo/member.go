package repo

import (
	"golang.org/x/net/context"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/database"
)

// MemberRepo 用户模块中操作数据库的接口（行为准则）
type MemberRepo interface {
	IsInMemberEmail(ctx context.Context, email string) (bool, error)     // Email是否已经在库中（已注册）
	IsInMemberAccount(ctx context.Context, account string) (bool, error) // 账号是否已经在库中（已注册）
	IsInMemberMobile(ctx context.Context, mobile string) (bool, error)   // 手机号是否已经在库中（已注册）
	SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account string, pwd string) (*member.Member, error)
	FindMemberById(ctx context.Context, id int64) (mem *member.Member, err error)
	FindMemberByIds(ctx context.Context, ids []int64) (list []*member.Member, err error)
}
