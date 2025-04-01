package rest

import (
	"context"

	"github.com/Hexes-rgb/employee-service/internal/domain"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, emp *domain.Employee) (int, error)
	GetEmployee(ctx context.Context, id int) (*domain.Employee, error)
	UpdateEmployee(ctx context.Context, emp *domain.Employee) error
	DeleteEmployee(ctx context.Context, id int) error
	GetCompanyEmployees(ctx context.Context, companyID int) ([]*domain.Employee, error)
	GetDepartmentEmployees(ctx context.Context, companyID int, deptName string) ([]*domain.Employee, error)
}

type DepartmentService interface {
	GetOrCreate(ctx context.Context, dept *domain.Department) (int, error)
	GetDepartment(ctx context.Context, id int) (*domain.Department, error)
}
