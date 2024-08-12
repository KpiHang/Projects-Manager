package repo

import (
	"context"
	"test.com/project-project/internal/data/pro"
)

type ProjectRepo interface {
	// FindProjectByMemId 返回projects, 项目总数，err
	FindProjectByMemId(ctx context.Context, memId int64, page int64, pageSize int64) ([]*pro.ProjectAndMember, int64, error)
}
