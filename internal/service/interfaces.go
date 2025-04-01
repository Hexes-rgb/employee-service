package service

import "github.com/Hexes-rgb/employee-service/internal/domain"

type EmployeeRepository interface {
	Create(emp *domain.Employee) (int, error)
	GetByID(id int) (*domain.Employee, error)
	Update(emp *domain.Employee) error
	Delete(id int) error
	GetByCompany(companyID int) ([]*domain.Employee, error)
	GetByDepartment(companyID, deptName int) ([]*domain.Employee, error)
}

type DepartmentRepository interface {
	GetOrCreate(dept *domain.Department) (int, error)
	GetByID(id int) (*domain.Department, error)
}
