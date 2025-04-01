package rest

import "github.com/Hexes-rgb/employee-service/internal/domain"

type EmployeeService interface {
	CreateEmployee(emp *domain.Employee) (int, error)
	GetEmployee(id int) (*domain.Employee, error)
	UpdateEmployee(emp *domain.Employee) error
	DeleteEmployee(id int) error
	GetCompanyEmployees(companyID int) ([]*domain.Employee, error)
	GetDepartmentEmployees(companyID, deptId int) ([]*domain.Employee, error)
}

type DepartmentService interface {
	GetOrCreate(dept *domain.Department) (int, error)
	GetDepartment(id int) (*domain.Department, error)
}
