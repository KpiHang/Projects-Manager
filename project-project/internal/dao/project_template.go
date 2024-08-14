package dao

import (
	"context"
	"test.com/project-project/internal/data/pro"
	"test.com/project-project/internal/database/gorms"
)

type ProjectTemplateDao struct {
	conn *gorms.GormConn
}

func NewProjectTemplateDao() *ProjectTemplateDao {
	return &ProjectTemplateDao{
		conn: gorms.NewGormConn(),
	}
}

func (p *ProjectTemplateDao) FindProjectTemplateSystem(ctx context.Context, page int64, size int64) (pts []pro.ProjectTemplate, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.
		Model(&pro.ProjectTemplate{}).
		Where("is_system = ?", 1).
		Limit(int(size)).
		Offset(int((page - 1) * size)).
		Find(&pts).Error
	if err != nil {
		return
	}
	err = session.Model(&pro.ProjectTemplate{}).Where("is_system = ?", 1).Count(&total).Error
	return
}

func (p *ProjectTemplateDao) FindProjectTemplateCustom(ctx context.Context, memId int64, organizationCode int64, page int64, size int64) (pts []pro.ProjectTemplate, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.
		Model(&pro.ProjectTemplate{}).
		Where("is_system = ? and member_code= ? and organization_code= ?", 0, memId, organizationCode). // 自定义模板，这个字段为0
		Limit(int(size)).
		Offset(int((page - 1) * size)).
		Find(&pts).Error
	if err != nil {
		return
	}
	err = session.Model(&pro.ProjectTemplate{}).Where("is_system = ? and member_code= ? and organization_code= ?", 0, memId, organizationCode).Count(&total).Error
	return
}

func (p *ProjectTemplateDao) FindProjectTemplateAll(ctx context.Context, organizationCode int64, page int64, size int64) (pts []pro.ProjectTemplate, total int64, err error) {
	session := p.conn.Session(ctx)
	err = session.
		Model(&pro.ProjectTemplate{}).
		Where("organization_code = ?", organizationCode).
		Limit(int(size)).
		Offset(int((page - 1) * size)).
		Find(&pts).Error
	if err != nil {
		return
	}
	err = session.Model(&pro.ProjectTemplate{}).Where("organization_code = ?", organizationCode).Count(&total).Error
	return
}
