package rest

import (
	"errors"

	"github.com/Hexes-rgb/employee-service/internal/domain"
)

func validateDepartment(dept *domain.Department) error {
	if dept.CompanyID == 0 {
		return errors.New("department companyId is required")
	}
	if dept.Name == "" {
		return errors.New("department name is required")
	}
	if dept.Phone == "" {
		return errors.New("department phone is required")
	}
	return nil
}

func validateEmployee(emp *domain.Employee) error {
	if emp.Name == "" {
		return errors.New("employee name is required")
	}
	if emp.Surname == "" {
		return errors.New("employee surname is required")
	}
	if emp.Phone == "" {
		return errors.New("employee phone is required")
	}
	if emp.CompanyID == 0 {
		return errors.New("employee companyId is required")
	}
	if emp.PassportNumber == "" {
		return errors.New("employee passport number is required")
	}
	return nil
}
