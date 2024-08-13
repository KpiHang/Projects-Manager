package dao

import (
	"context"
	"test.com/project-project/internal/data/pro"
	"test.com/project-project/internal/database/gorms"
)

type ProjectDao struct {
	conn *gorms.GormConn
}

func NewProjectDao() *ProjectDao {
	return &ProjectDao{
		conn: gorms.NewGormConn(),
	}
}

// FindProjectByMemId 返回projects, 项目总数，err
func (p *ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error) {

	var pms []*pro.ProjectAndMember
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
	raw := session.Raw("select * from ms_project  a, ms_project_member b where a.id = b.project_code and b.member_code = ? order by sort limit ?,?", memId, index, pageSize)
	raw.Scan(&pms)
	var total int64
	err := session.Model(&pro.ProjectMember{}).Where("member_code=?", memId).Count(&total).Error
	return pms, total, err
}
