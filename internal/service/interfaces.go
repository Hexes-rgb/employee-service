package service

import (
	"context"

	"github.com/Hexes-rgb/employee-service/internal/domain"
)

type EmployeeRepository interface {
	Create(ctx context.Context, emp *domain.Employee) (int, error)
	GetByID(ctx context.Context, id int) (*domain.Employee, error)
	Update(ctx context.Context, emp *domain.Employee) error
	Delete(ctx context.Context, id int) error
	GetByCompany(ctx context.Context, companyID int) ([]*domain.Employee, error)
	GetByDepartment(ctx context.Context, companyID int, deptName string) ([]*domain.Employee, error)
}

type DepartmentRepository interface {
	GetOrCreate(ctx context.Context, dept *domain.Department) (int, error)
	GetByID(ctx context.Context, id int) (*domain.Department, error)
}
