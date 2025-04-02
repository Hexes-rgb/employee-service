package service_test

import (
	"errors"
	"testing"

	"github.com/Hexes-rgb/employee-service/internal/domain"
	"github.com/Hexes-rgb/employee-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type EmployeeRepositoryMock struct {
	mock.Mock
}

func (m *EmployeeRepositoryMock) Create(emp *domain.Employee) (int, error) {
	args := m.Called(emp)
	return args.Int(0), args.Error(1)
}

func (m *EmployeeRepositoryMock) GetByID(id int) (*domain.Employee, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *EmployeeRepositoryMock) Update(emp *domain.Employee) error {
	args := m.Called(emp)
	return args.Error(0)
}

func (m *EmployeeRepositoryMock) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *EmployeeRepositoryMock) GetByCompany(companyID int) ([]*domain.Employee, error) {
	args := m.Called(companyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Employee), args.Error(1)
}

func (m *EmployeeRepositoryMock) GetByDepartment(companyID, deptID int) ([]*domain.Employee, error) {
	args := m.Called(companyID, deptID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Employee), args.Error(1)
}

func ptrInt(i int) *int {
	return &i
}

func TestEmployeeService_CreateEmployee(t *testing.T) {
	t.Run("Success: Creating an employee with a department", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		inputDept := &domain.Department{
			CompanyID: 1,
			Name:      "Engineering",
			Phone:     "+123456789",
		}

		expectedEmployee := &domain.Employee{
			Name:           "John",
			Surname:        "Doe",
			Phone:          "+79998887766",
			CompanyID:      1,
			PassportType:   "internal",
			PassportNumber: "1234567890",
			Department:     inputDept,
		}

		deptRepo.On("GetOrCreate", inputDept).Return(42, nil)

		empRepo.On("Create", mock.AnythingOfType("*domain.Employee")).Return(100, nil).Run(func(args mock.Arguments) {
			emp := args.Get(0).(*domain.Employee)
			assert.Equal(t, "John", emp.Name)
			assert.Equal(t, "Doe", emp.Surname)
			assert.Equal(t, "+79998887766", emp.Phone)
			assert.Equal(t, 1, emp.CompanyID)
			assert.Equal(t, "internal", emp.PassportType)
			assert.Equal(t, "1234567890", emp.PassportNumber)
			assert.Equal(t, 42, *emp.DepartmentID)
			assert.Equal(t, inputDept, emp.Department)
		})

		svc := service.NewEmployeeService(empRepo, deptRepo)
		id, err := svc.CreateEmployee(expectedEmployee)

		assert.NoError(t, err)
		assert.Equal(t, 100, id)
		empRepo.AssertExpectations(t)
		deptRepo.AssertExpectations(t)
	})

	t.Run("Success: Creating an employee without a department", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		inputEmployee := &domain.Employee{
			Name:      "Jane",
			Surname:   "Smith",
			CompanyID: 1,
		}

		empRepo.On("Create", inputEmployee).Return(101, nil)

		svc := service.NewEmployeeService(empRepo, deptRepo)
		id, err := svc.CreateEmployee(inputEmployee)

		assert.NoError(t, err)
		assert.Equal(t, 101, id)
		empRepo.AssertExpectations(t)
	})

	t.Run("Error: Failed to create department", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		inputDept := &domain.Department{
			CompanyID: 1,
			Name:      "HR",
			Phone:     "+987654321",
		}

		inputEmployee := &domain.Employee{
			Name:       "Alice",
			Surname:    "Johnson",
			CompanyID:  1,
			Department: inputDept,
		}

		deptRepo.On("GetOrCreate", inputDept).Return(0, errors.New("db error"))

		svc := service.NewEmployeeService(empRepo, deptRepo)
		_, err := svc.CreateEmployee(inputEmployee)

		assert.EqualError(t, err, "failed to get or create department: db error")
		deptRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_GetEmployee(t *testing.T) {
	t.Run("Success: Getting an employee with a department", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		empRepo.On("GetByID", 1).Return(&domain.Employee{
			ID:           1,
			Name:         "John",
			Surname:      "Doe",
			Phone:        "+79998887766",
			CompanyID:    1,
			DepartmentID: ptrInt(42),
		}, nil)

		deptRepo.On("GetByID", 42).Return(&domain.Department{
			ID:        42,
			CompanyID: 1,
			Name:      "Engineering",
			Phone:     "+123456789",
		}, nil)

		svc := service.NewEmployeeService(empRepo, deptRepo)
		emp, err := svc.GetEmployee(1)

		assert.NoError(t, err)
		assert.Equal(t, &domain.Employee{
			ID:           1,
			Name:         "John",
			Surname:      "Doe",
			Phone:        "+79998887766",
			CompanyID:    1,
			DepartmentID: ptrInt(42),
			Department: &domain.Department{
				ID:        42,
				CompanyID: 1,
				Name:      "Engineering",
				Phone:     "+123456789",
			},
		}, emp)
		empRepo.AssertExpectations(t)
		deptRepo.AssertExpectations(t)
	})

	t.Run("Error: employee not found", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		empRepo.On("GetByID", 999).Return(nil, errors.New("not found"))

		svc := service.NewEmployeeService(empRepo, deptRepo)
		_, err := svc.GetEmployee(999)

		assert.EqualError(t, err, "failed to get employee: not found")
		empRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_UpdateEmployee(t *testing.T) {
	t.Run("Success: Employee update with department", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		newDept := &domain.Department{
			CompanyID: 1,
			Name:      "New Department",
			Phone:     "+987654321",
		}

		deptRepo.On("GetOrCreate", newDept).Return(43, nil)
		empRepo.On("Update", mock.AnythingOfType("*domain.Employee")).Return(nil).Run(func(args mock.Arguments) {
			emp := args.Get(0).(*domain.Employee)
			assert.Equal(t, 1, emp.ID)
			assert.Equal(t, 43, *emp.DepartmentID)
			assert.Equal(t, newDept, emp.Department)
		})

		svc := service.NewEmployeeService(empRepo, deptRepo)
		err := svc.UpdateEmployee(&domain.Employee{
			ID:         1,
			Department: newDept,
		})

		assert.NoError(t, err)
		empRepo.AssertExpectations(t)
		deptRepo.AssertExpectations(t)
	})

	t.Run("Error: Failed to update department", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		newDept := &domain.Department{
			CompanyID: 1,
			Name:      "Finance",
			Phone:     "+1122334455",
		}

		deptRepo.On("GetOrCreate", newDept).Return(0, errors.New("db error"))

		svc := service.NewEmployeeService(empRepo, deptRepo)
		err := svc.UpdateEmployee(&domain.Employee{
			ID:         1,
			Department: newDept,
		})

		assert.EqualError(t, err, "failed to get or create department: db error")
		deptRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_DeleteEmployee(t *testing.T) {
	t.Run("Success: Employee Removal", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		empRepo.On("Delete", 1).Return(nil)

		svc := service.NewEmployeeService(empRepo, deptRepo)
		err := svc.DeleteEmployee(1)

		assert.NoError(t, err)
		empRepo.AssertExpectations(t)
	})

	t.Run("Error: Failed to delete employee", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		empRepo.On("Delete", 999).Return(errors.New("db error"))

		svc := service.NewEmployeeService(empRepo, deptRepo)
		err := svc.DeleteEmployee(999)

		assert.EqualError(t, err, "failed to delete employee: db error")
		empRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_GetCompanyEmployees(t *testing.T) {
	t.Run("Success: Getting employees of a company with departments", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		empRepo.On("GetByCompany", 1).Return([]*domain.Employee{
			{ID: 1, Name: "John", CompanyID: 1, DepartmentID: ptrInt(42)},
			{ID: 2, Name: "Jane", CompanyID: 1},
		}, nil)

		deptRepo.On("GetByID", 42).Return(&domain.Department{
			ID:        42,
			CompanyID: 1,
			Name:      "Engineering",
		}, nil)

		svc := service.NewEmployeeService(empRepo, deptRepo)
		employees, err := svc.GetCompanyEmployees(1)

		assert.NoError(t, err)
		assert.Equal(t, []*domain.Employee{
			{
				ID:           1,
				Name:         "John",
				CompanyID:    1,
				DepartmentID: ptrInt(42),
				Department: &domain.Department{
					ID:        42,
					CompanyID: 1,
					Name:      "Engineering",
				},
			},
			{
				ID:        2,
				Name:      "Jane",
				CompanyID: 1,
			},
		}, employees)
		empRepo.AssertExpectations(t)
		deptRepo.AssertExpectations(t)
	})
}

func TestEmployeeService_GetDepartmentEmployees(t *testing.T) {
	t.Run("Success: Getting Department Staff", func(t *testing.T) {
		empRepo := new(EmployeeRepositoryMock)
		deptRepo := new(DepartmentRepositoryMock)

		empRepo.On("GetByDepartment", 1, 42).Return([]*domain.Employee{
			{ID: 1, Name: "John", CompanyID: 1, DepartmentID: ptrInt(42)},
		}, nil)

		deptRepo.On("GetByID", 42).Return(&domain.Department{
			ID:        42,
			CompanyID: 1,
			Name:      "Engineering",
		}, nil)

		svc := service.NewEmployeeService(empRepo, deptRepo)
		employees, err := svc.GetDepartmentEmployees(1, 42)

		assert.NoError(t, err)
		assert.Equal(t, []*domain.Employee{
			{
				ID:           1,
				Name:         "John",
				CompanyID:    1,
				DepartmentID: ptrInt(42),
				Department: &domain.Department{
					ID:        42,
					CompanyID: 1,
					Name:      "Engineering",
				},
			},
		}, employees)
		empRepo.AssertExpectations(t)
		deptRepo.AssertExpectations(t)
	})
}
