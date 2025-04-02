package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Hexes-rgb/employee-service/internal/domain"
	"github.com/lib/pq"
)

type EmployeeRepo struct {
	db *sql.DB
}

func NewEmployeeRepo(db *sql.DB) *EmployeeRepo {
	return &EmployeeRepo{db: db}
}

func (r *EmployeeRepo) Create(emp *domain.Employee) (int, error) {
	var id int
	query := `INSERT INTO employees 
        (name, surname, phone, company_id, department_id, passport_type, passport_number)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := r.db.QueryRow(query,
		emp.Name,
		emp.Surname,
		emp.Phone,
		emp.CompanyID,
		emp.DepartmentID,
		emp.PassportType,
		emp.PassportNumber,
	).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "employees_phone_key":
				return 0, fmt.Errorf("employee with this phone number already exists")
			case "employees_passport_number_key":
				return 0, fmt.Errorf("employee with this passport number already exists")
			}
		}
		return 0, fmt.Errorf("failed to create employee: %w", err)
	}

	return id, nil
}

func (r *EmployeeRepo) GetByID(id int) (*domain.Employee, error) {
	query := `SELECT id, name, surname, phone, company_id, department_id, 
        passport_type, passport_number FROM employees WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var emp domain.Employee
	err := row.Scan(
		&emp.ID,
		&emp.Name,
		&emp.Surname,
		&emp.Phone,
		&emp.CompanyID,
		&emp.DepartmentID,
		&emp.PassportType,
		&emp.PassportNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return &emp, nil
}

func (r *EmployeeRepo) Update(emp *domain.Employee) error {
	var updates []string
	var args []interface{}
	argID := 1

	if emp.Name != "" {
		updates = append(updates, "name = $"+strconv.Itoa(argID))
		args = append(args, emp.Name)
		argID++
	}
	if emp.Surname != "" {
		updates = append(updates, "surname = $"+strconv.Itoa(argID))
		args = append(args, emp.Surname)
		argID++
	}
	if emp.Phone != "" {
		updates = append(updates, "phone = $"+strconv.Itoa(argID))
		args = append(args, emp.Phone)
		argID++
	}
	if emp.CompanyID != 0 {
		updates = append(updates, "company_id = $"+strconv.Itoa(argID))
		args = append(args, emp.CompanyID)
		argID++
	}
	if emp.DepartmentID != nil {
		updates = append(updates, "department_id = $"+strconv.Itoa(argID))
		args = append(args, emp.DepartmentID)
		argID++
	}
	if emp.PassportType != "" {
		updates = append(updates, "passport_type = $"+strconv.Itoa(argID))
		args = append(args, emp.PassportType)
		argID++
	}
	if emp.PassportNumber != "" {
		updates = append(updates, "passport_number = $"+strconv.Itoa(argID))
		args = append(args, emp.PassportNumber)
		argID++
	}

	args = append(args, emp.ID)

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := "UPDATE employees SET " + strings.Join(updates, ", ") + " WHERE id = $" + strconv.Itoa(argID)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "employees_phone_key":
				return fmt.Errorf("employee with this phone number already exists")
			case "employees_passport_number_key":
				return fmt.Errorf("employee with this passport number already exists")
			}
		}
		return fmt.Errorf("failed to update employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

func (r *EmployeeRepo) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM employees WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

func (r *EmployeeRepo) GetByCompany(companyID int) ([]*domain.Employee, error) {
	query := `SELECT id, name, surname, phone, company_id, department_id, 
        passport_type, passport_number FROM employees WHERE company_id = $1`

	rows, err := r.db.Query(query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	defer rows.Close()

	var employees []*domain.Employee
	for rows.Next() {
		var emp domain.Employee
		err := rows.Scan(
			&emp.ID,
			&emp.Name,
			&emp.Surname,
			&emp.Phone,
			&emp.CompanyID,
			&emp.DepartmentID,
			&emp.PassportType,
			&emp.PassportNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, &emp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(employees) == 0 {
		return nil, fmt.Errorf("employees not found for company id %d", companyID)
	}

	return employees, nil
}

func (r *EmployeeRepo) GetByDepartment(companyID, deptId int) ([]*domain.Employee, error) {
	query := `SELECT e.id, e.name, e.surname, e.phone, e.company_id, 
        e.department_id, e.passport_type, e.passport_number
        FROM employees e
        JOIN departments d ON e.department_id = d.id
        WHERE e.company_id = $1 AND d.id = $2`

	rows, err := r.db.Query(query, companyID, deptId)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	defer rows.Close()

	var employees []*domain.Employee
	for rows.Next() {
		var emp domain.Employee
		err := rows.Scan(
			&emp.ID,
			&emp.Name,
			&emp.Surname,
			&emp.Phone,
			&emp.CompanyID,
			&emp.DepartmentID,
			&emp.PassportType,
			&emp.PassportNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, &emp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if len(employees) == 0 {
		return nil, fmt.Errorf("employees not found for company id %d and department id %d", companyID, deptId)
	}

	return employees, nil
}
