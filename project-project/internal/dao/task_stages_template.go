package dao

import (
	"context"
	"test.com/project-project/internal/data/task"
	"test.com/project-project/internal/database/gorms"
)

type TaskStagesTemplateDao struct {
	conn *gorms.GormConn
}

// FindInProTemIds 根据项目模板id查找模板包含的stage
func (t *TaskStagesTemplateDao) FindInProTemIds(ctx context.Context, ids []int) ([]task.MsTaskStagesTemplate, error) {
	var tsts []task.MsTaskStagesTemplate
	session := t.conn.Session(ctx)
	err := session.Where("project_template_code in ?", ids).Find(&tsts).Error
	return tsts, err
}

func NewTaskStagesTemplateDao() *TaskStagesTemplateDao {
	return &TaskStagesTemplateDao{
		conn: gorms.NewGormConn(),
	}
}
