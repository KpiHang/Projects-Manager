package repo

import (
	"context"
	"test.com/project-project/internal/data/task"
)

type TaskStagesTemplateRepo interface {
	// FindInProTemIds 根据项目模板id查找模板包含的stage
	FindInProTemIds(ctx context.Context, ids []int) ([]task.MsTaskStagesTemplate, error)
}
