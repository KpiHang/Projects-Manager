package dao

import (
	"context"
	"test.com/project-project/internal/data/menu"
	"test.com/project-project/internal/database/gorms"
)

type MenuDao struct {
	conn *gorms.GormConn
}

func NewMenuDao() *MenuDao {
	return &MenuDao{
		conn: gorms.NewGormConn(),
	}
}

func (m *MenuDao) FindMenus(ctx context.Context) (pms []*menu.ProjectMenu, err error) {
	session := m.conn.Session(ctx)
	err = session.Find(&pms).Error
	return
}
