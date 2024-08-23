package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type ProjectLogDao struct {
	conn *gorms.GormConn
}

func (p *ProjectLogDao) FindLogByMemberCode(ctx context.Context, id int64, page int64, size int64) (list []*data.ProjectLog, total int64, err error) {
	session := p.conn.Session(ctx)
	offset := (page - 1) * size
	err = session.Model(&data.ProjectLog{}).
		Where("member_code=?", id).
		Limit(int(size)).Offset(int(offset)).
		Order("create_time desc").
		Find(&list).Error
	err = session.Model(&data.ProjectLog{}).
		Where("member_code=?", id).Count(&total).Error
	return
}

func (p *ProjectLogDao) SaveProjectLog(pl *data.ProjectLog) {
	session := p.conn.Session(context.Background())
	session.Save(&pl)
}

// FindLogByTaskCode 获取任务日志不分页
func (p *ProjectLogDao) FindLogByTaskCode(ctx context.Context, taskCode int64, comment int) (list []*data.ProjectLog, total int64, err error) {
	session := p.conn.Session(ctx)
	model := session.Model(&data.ProjectLog{})
	if comment == 1 {
		model.Where("source_code=? and is_comment=?", taskCode, comment).Find(&list)
		model.Where("source_code=? and is_comment=?", taskCode, comment).Count(&total)
	} else {
		model.Where("source_code=?", taskCode).Find(&list)
		model.Where("source_code=?", taskCode).Count(&total)
	}
	return
}

// FindLogByTaskCodePage 获取任务日志分页
func (p *ProjectLogDao) FindLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) (list []*data.ProjectLog, total int64, err error) {
	session := p.conn.Session(ctx)
	model := session.Model(&data.ProjectLog{})
	offset := (page - 1) * pageSize
	if comment == 1 {
		model.Where("source_code=? and is_comment=?", taskCode, comment).Limit(pageSize).Offset(offset).Find(&list)
		model.Where("source_code=? and is_comment=?", taskCode, comment).Count(&total)
	} else {
		model.Where("source_code=?", taskCode).Limit(pageSize).Offset(offset).Find(&list)
		model.Where("source_code=?", taskCode).Count(&total)
	}
	return
}

func NewProjectLogDao() *ProjectLogDao {
	return &ProjectLogDao{
		conn: gorms.NewGormConn(),
	}
}
