package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/database"
	"test.com/project-user/internal/database/gorms"
)

type MemberDao struct {
	conn *gorms.GormConn
}

func NewMemberDao() *MemberDao {
	return &MemberDao{
		conn: gorms.NewGormConn(),
	}
}

func (m *MemberDao) IsInMemberEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) IsInMemberAccount(ctx context.Context, account string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("account = ?", account).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) IsInMemberMobile(ctx context.Context, mobile string) (bool, error) {
	var count int64
	err := m.conn.Session(ctx).Model(&member.Member{}).Where("mobile = ?", mobile).Count(&count).Error
	return count > 0, err
}

func (m *MemberDao) SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error {
	m.conn = conn.(*gorms.GormConn)
	return m.conn.Tx(ctx).Create(mem).Error
}

func (m *MemberDao) FindMember(ctx context.Context, account string, pwd string) (*member.Member, error) {
	mem := &member.Member{} // 给Find的值必须已经分配内存空间了；
	err := m.conn.Session(ctx).Where("account = ? AND password = ?", account, pwd).Find(mem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // go语言中 目前 redis mysql 查询不到记录会报这个错，但这个错是非业务错；
		return nil, nil
	}
	return mem, err
}

func (m *MemberDao) FindMemberById(ctx context.Context, id int64) (mem *member.Member, err error) {
	err = m.conn.Session(ctx).Where("id = ?", id).Find(&mem).Error
	return
}
