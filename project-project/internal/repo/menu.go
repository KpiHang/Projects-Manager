package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type MenuRepo interface {
	FindMenus(ctx context.Context) ([]*data.ProjectMenu, error) // find所有菜单，也包括子菜单
}
