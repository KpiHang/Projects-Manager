package domain

import (
	"context"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/repo"
	"time"
)

type DepartmentDomain struct {
	departmentRepo repo.DepartmentRepo
}

func (d *DepartmentDomain) FindDepartmentById(id int64) (*data.Department, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return d.departmentRepo.FindDepartmentById(c, id)
}

func NewDepartmentDomain() *DepartmentDomain {
	return &DepartmentDomain{
		departmentRepo: dao.NewDepartmentDao(),
	}
}
