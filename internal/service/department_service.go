package service

import (
	"fmt"

	"github.com/Hexes-rgb/employee-service/internal/domain"
)

type DepartmentService struct {
	repo DepartmentRepository
}

func NewDepartmentService(repo DepartmentRepository) *DepartmentService {
	return &DepartmentService{repo: repo}
}

func (s *DepartmentService) GetOrCreate(dept *domain.Department) (int, error) {
	id, err := s.repo.GetOrCreate(dept)
	if err != nil {
		return 0, fmt.Errorf("failed to get or create department: %w", err)
	}
	return id, nil
}

func (s *DepartmentService) GetDepartment(id int) (*domain.Department, error) {
	dept, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	return dept, nil
}
