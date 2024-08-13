package dao

import (
	"context"
	"fmt"
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
func (p *ProjectDao) FindProjectByMemId(ctx context.Context, memId int64, condition string, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error) {
	var pms []*pro.ProjectAndMember
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
	sql := fmt.Sprintf("select * from ms_project  a, ms_project_member b where a.id = b.project_code and b.member_code = ? %s order by sort limit ?,?", condition)
	raw := session.Raw(sql, memId, index, pageSize)
	raw.Scan(&pms)

	var total int64
	query := fmt.Sprintf("select count(*) from ms_project  a, ms_project_member b where a.id = b.project_code and b.member_code = ? %s", condition)
	tx := session.Raw(query, memId)
	err := tx.Scan(&total).Error
	return pms, total, err
}

func (p *ProjectDao) FindCollectProjectByMemId(ctx context.Context, memId int64, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error) {
	var pms []*pro.ProjectAndMember
	session := p.conn.Session(ctx)
	index := (page - 1) * pageSize
	sql := fmt.Sprintf("select * from ms_project  where id in (select project_code from ms_project_collection where member_code = ?) order by sort limit ?,?")
	raw := session.Raw(sql, memId, index, pageSize)
	raw.Scan(&pms)

	var total int64                       // 收藏项目总数，是在收藏表中进行的；
	query := fmt.Sprintf("member_code=?") // pro.ProjectCollection 实现了TableName接口，所以不需要手动指定表名；
	err := session.Model(&pro.ProjectCollection{}).Where(query, memId).Count(&total).Error
	return pms, total, err
}
