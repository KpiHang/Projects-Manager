package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type ProjectLogRepo interface {
	FindLogByTaskCode(ctx context.Context, taskCode int64, comment int) (list []*data.ProjectLog, total int64, err error)
	FindLogByTaskCodePage(ctx context.Context, taskCode int64, comment int, page int, pageSize int) (list []*data.ProjectLog, total int64, err error)
	SaveProjectLog(pl *data.ProjectLog)
	FindLogByMemberCode(ctx context.Context, id int64, page int64, size int64) (list []*data.ProjectLog, total int64, err error)
}
