package postgres

import (
	"database/sql"
	"fmt"

	"github.com/Hexes-rgb/employee-service/internal/domain"
	"github.com/lib/pq"
)

type DepartmentRepo struct {
	db *sql.DB
}

func NewDepartmentRepo(db *sql.DB) *DepartmentRepo {
	return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) GetOrCreate(dept *domain.Department) (int, error) {
	var id int
	err := r.db.QueryRow(
		"SELECT id FROM departments WHERE company_id = $1 AND name = $2",
		dept.CompanyID, dept.Name,
	).Scan(&id)

	if err == nil {
		return id, nil
	}

	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query department: %w", err)
	}

	err = r.db.QueryRow(
		"INSERT INTO departments (company_id, name, phone) VALUES ($1, $2, $3) RETURNING id",
		dept.CompanyID, dept.Name, dept.Phone,
	).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "departments_company_id_name_key":
				return 0, fmt.Errorf("department with this name already exists in this company")
			case "departments_phone_key":
				return 0, fmt.Errorf("department with this phone number already exists")
			}
		}
		return 0, fmt.Errorf("failed to create department: %w", err)
	}

	return id, nil
}

func (r *DepartmentRepo) GetByID(id int) (*domain.Department, error) {
	query := "SELECT id, company_id, name, phone FROM departments WHERE id = $1"

	row := r.db.QueryRow(query, id)

	var dept domain.Department
	err := row.Scan(
		&dept.ID,
		&dept.CompanyID,
		&dept.Name,
		&dept.Phone,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("department not found")
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return &dept, nil
}
