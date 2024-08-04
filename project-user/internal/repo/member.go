package repo

import (
	"golang.org/x/net/context"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/database"
)

type MemberRepo interface {
	IsInMemberEmail(ctx context.Context, email string) (bool, error)     // Email是否已经在库中（已注册）
	IsInMemberAccount(ctx context.Context, account string) (bool, error) // 账号是否已经在库中（已注册）
	IsInMemberMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error // 手机号是否已经在库中（已注册）
}
