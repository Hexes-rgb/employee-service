package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Hexes-rgb/employee-service/internal/domain"
)

type EmployeeHandlers struct {
	service EmployeeService
}

func NewEmployeeHandlers(s EmployeeService) *EmployeeHandlers {
	return &EmployeeHandlers{service: s}
}

func (h *EmployeeHandlers) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id, err := h.service.CreateEmployee(&emp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *EmployeeHandlers) GetEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	emp, err := h.service.GetEmployee(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, emp)
}

func (h *EmployeeHandlers) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	var emp domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	emp.ID = id

	if err := h.service.UpdateEmployee(&emp); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Employee updated successfully"})
}

func (h *EmployeeHandlers) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid employee ID")
		return
	}

	if err := h.service.DeleteEmployee(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Employee deleted successfully"})
}

func (h *EmployeeHandlers) GetCompanyEmployees(w http.ResponseWriter, r *http.Request) {
	companyID, err := strconv.Atoi(r.PathValue("companyId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	employees, err := h.service.GetCompanyEmployees(companyID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, employees)
}

func (h *EmployeeHandlers) GetDepartmentEmployees(w http.ResponseWriter, r *http.Request) {
	companyID, err := strconv.Atoi(r.PathValue("companyId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	deptId, err := strconv.Atoi(r.PathValue("departmentId"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	employees, err := h.service.GetDepartmentEmployees(companyID, deptId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, employees)
}
