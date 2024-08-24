package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type DepartmentDao struct {
	conn *gorms.GormConn
}

func (d *DepartmentDao) FindDepartmentById(ctx context.Context, id int64) (dt *data.Department, err error) {
	session := d.conn.Session(ctx)
	err = session.Where("id=?", id).Find(&dt).Error
	return
}

func NewDepartmentDao() *DepartmentDao {
	return &DepartmentDao{
		conn: gorms.NewGormConn(),
	}
}
