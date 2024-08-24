package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type DepartmentRepo interface {
	FindDepartmentById(ctx context.Context, id int64) (*data.Department, error)
}
