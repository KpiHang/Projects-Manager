package dao

import (
	"context"
	"gorm.io/gorm"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type TaskDao struct {
	conn *gorms.GormConn
}

func NewTaskDao() *TaskDao {
	return &TaskDao{
		conn: gorms.NewGormConn(),
	}
}

func (t *TaskDao) FindTaskByStageCode(ctx context.Context, stageCode int) (list []*data.Task, err error) {
	session := t.conn.Session(ctx)
	err = session.Model(&data.Task{}).
		Where("stage_code = ? and deleted= 0", stageCode).
		Order("sort asc").
		Find(&list).Error
	return
}

func (t *TaskDao) FindTaskMemberByTaskId(ctx context.Context, taskCode int64, memberId int64) (task *data.TaskMember, err error) {
	err = t.conn.Session(ctx).
		Where("task_code=? and member_code=?", taskCode, memberId).
		Limit(1).
		Find(&task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (t *TaskDao) FindTaskMaxIdNum(ctx context.Context, projectCode int64) (v *int, err error) {
	session := t.conn.Session(ctx)
	err = session.Model(&data.Task{}).
		Where("project_code=?", projectCode).
		Select("max(id_num)").Scan(&v).Error
	return
}

func (t *TaskDao) FindTaskSort(ctx context.Context, projectCode int64, stageCode int64) (v *int, err error) {
	session := t.conn.Session(ctx)
	err = session.Model(&data.Task{}).
		Where("project_code=? and stage_code=?", projectCode, stageCode).
		Select("max(sort)").Scan(&v).Error
	return
}

func (t *TaskDao) SaveTask(ctx context.Context, conn database.DbConn, ts *data.Task) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Save(&ts).Error
	return err
}

func (t *TaskDao) SaveTaskMember(ctx context.Context, conn database.DbConn, tm *data.TaskMember) error {
	t.conn = conn.(*gorms.GormConn)
	err := t.conn.Tx(ctx).Save(&tm).Error
	return err
}
