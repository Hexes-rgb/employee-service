package service

import (
	"fmt"

	"github.com/Hexes-rgb/employee-service/internal/domain"
)

type EmployeeService struct {
	empRepo  EmployeeRepository
	deptRepo DepartmentRepository
}

func NewEmployeeService(
	empRepo EmployeeRepository,
	deptRepo DepartmentRepository,
) *EmployeeService {
	return &EmployeeService{
		empRepo:  empRepo,
		deptRepo: deptRepo,
	}
}

func (s *EmployeeService) CreateEmployee(emp *domain.Employee) (int, error) {
	if emp.Department != nil {
		deptID, err := s.deptRepo.GetOrCreate(emp.Department)
		if err != nil {
			return 0, fmt.Errorf("failed to get or create department: %w", err)
		}
		emp.DepartmentID = &deptID
	}

	id, err := s.empRepo.Create(emp)
	if err != nil {
		return 0, fmt.Errorf("failed to create employee: %w", err)
	}

	return id, nil
}

func (s *EmployeeService) GetEmployee(id int) (*domain.Employee, error) {
	emp, err := s.empRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	if emp.DepartmentID != nil {
		dept, err := s.deptRepo.GetByID(*emp.DepartmentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get department: %w", err)
		}
		emp.Department = dept
	}

	return emp, nil
}

func (s *EmployeeService) UpdateEmployee(emp *domain.Employee) error {
	if emp.Department != nil {
		deptID, err := s.deptRepo.GetOrCreate(emp.Department)
		if err != nil {
			return fmt.Errorf("failed to get or create department: %w", err)
		}
		emp.DepartmentID = &deptID
	}

	if err := s.empRepo.Update(emp); err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

func (s *EmployeeService) DeleteEmployee(id int) error {
	if err := s.empRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}
	return nil
}

func (s *EmployeeService) GetCompanyEmployees(companyID int) ([]*domain.Employee, error) {
	employees, err := s.empRepo.GetByCompany(companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	for _, emp := range employees {
		if emp.DepartmentID != nil {
			dept, err := s.deptRepo.GetByID(*emp.DepartmentID)
			if err != nil {
				return nil, fmt.Errorf("failed to get department: %w", err)
			}
			emp.Department = dept
		}
	}

	return employees, nil
}

func (s *EmployeeService) GetDepartmentEmployees(companyID, deptId int) ([]*domain.Employee, error) {
	employees, err := s.empRepo.GetByDepartment(companyID, deptId)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	for _, emp := range employees {
		if emp.DepartmentID != nil {
			dept, err := s.deptRepo.GetByID(*emp.DepartmentID)
			if err != nil {
				return nil, fmt.Errorf("failed to get department: %w", err)
			}
			emp.Department = dept
		}
	}

	return employees, nil
}
