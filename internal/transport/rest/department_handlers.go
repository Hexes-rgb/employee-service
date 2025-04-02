package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Hexes-rgb/employee-service/internal/domain"
)

type DepartmentHandlers struct {
	service DepartmentService
}

func NewDepartmentHandlers(s DepartmentService) *DepartmentHandlers {
	return &DepartmentHandlers{service: s}
}

func (h *DepartmentHandlers) GetOrCreateDepartment(w http.ResponseWriter, r *http.Request) {
	var dept domain.Department
	if err := json.NewDecoder(r.Body).Decode(&dept); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id, err := h.service.GetOrCreate(&dept)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, IDResponse{ID: id})
}

func (h *DepartmentHandlers) GetDepartment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid department ID")
		return
	}

	dept, err := h.service.GetDepartment(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, dept)
}
