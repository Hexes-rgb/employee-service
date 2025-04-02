package service_test

import (
	"errors"
	"testing"

	"github.com/Hexes-rgb/employee-service/internal/domain"
	"github.com/Hexes-rgb/employee-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DepartmentRepositoryMock struct {
	mock.Mock
}

func (m *DepartmentRepositoryMock) GetOrCreate(dept *domain.Department) (int, error) {
	args := m.Called(dept)
	return args.Int(0), args.Error(1)
}

func (m *DepartmentRepositoryMock) GetByID(id int) (*domain.Department, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Department), args.Error(1)
}

func TestDepartmentService_GetOrCreate(t *testing.T) {
	tests := []struct {
		name        string
		inputDept   *domain.Department
		mockSetup   func(*DepartmentRepositoryMock)
		expectedID  int
		expectedErr string
	}{
		{
			name: "Success: create new department with all fields",
			inputDept: &domain.Department{
				CompanyID: 1,
				Name:      "Engineering",
				Phone:     "+123456789",
			},
			mockSetup: func(m *DepartmentRepositoryMock) {
				m.On("GetOrCreate", &domain.Department{
					CompanyID: 1,
					Name:      "Engineering",
					Phone:     "+123456789",
				}).Return(42, nil)
			},
			expectedID: 42,
		},
		{
			name: "Success: create department with required fields only",
			inputDept: &domain.Department{
				CompanyID: 1,
				Name:      "HR",
			},
			mockSetup: func(m *DepartmentRepositoryMock) {
				m.On("GetOrCreate", &domain.Department{
					CompanyID: 1,
					Name:      "HR",
				}).Return(43, nil)
			},
			expectedID: 43,
		},
		{
			name: "Error: repository fails",
			inputDept: &domain.Department{
				CompanyID: 1,
				Name:      "Finance",
			},
			mockSetup: func(m *DepartmentRepositoryMock) {
				m.On("GetOrCreate", &domain.Department{
					CompanyID: 1,
					Name:      "Finance",
				}).Return(0, errors.New("database connection failed"))
			},
			expectedErr: "failed to get or create department: database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(DepartmentRepositoryMock)
			tt.mockSetup(repo)

			svc := service.NewDepartmentService(repo)
			id, err := svc.GetOrCreate(tt.inputDept)

			assert.Equal(t, tt.expectedID, id)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestDepartmentService_GetDepartment(t *testing.T) {
	tests := []struct {
		name        string
		inputID     int
		mockSetup   func(*DepartmentRepositoryMock)
		expected    *domain.Department
		expectedErr string
	}{
		{
			name:    "Success: get full department",
			inputID: 1,
			mockSetup: func(m *DepartmentRepositoryMock) {
				m.On("GetByID", 1).Return(&domain.Department{
					ID:        1,
					CompanyID: 1,
					Name:      "Engineering",
					Phone:     "+123456789",
				}, nil)
			},
			expected: &domain.Department{
				ID:        1,
				CompanyID: 1,
				Name:      "Engineering",
				Phone:     "+123456789",
			},
		},
		{
			name:    "Success: get department without phone",
			inputID: 2,
			mockSetup: func(m *DepartmentRepositoryMock) {
				m.On("GetByID", 2).Return(&domain.Department{
					ID:        2,
					CompanyID: 1,
					Name:      "HR",
				}, nil)
			},
			expected: &domain.Department{
				ID:        2,
				CompanyID: 1,
				Name:      "HR",
			},
		},
		{
			name:    "Error: department not found",
			inputID: 999,
			mockSetup: func(m *DepartmentRepositoryMock) {
				m.On("GetByID", 999).Return(nil, errors.New("not found"))
			},
			expectedErr: "failed to get department: not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(DepartmentRepositoryMock)
			tt.mockSetup(repo)

			svc := service.NewDepartmentService(repo)
			dept, err := svc.GetDepartment(tt.inputID)

			assert.Equal(t, tt.expected, dept)
			if tt.expectedErr != "" {
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}
