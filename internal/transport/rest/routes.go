package rest

import "net/http"

func NewRouter(
	empService EmployeeService,
	deptService DepartmentService,
) *http.ServeMux {
	router := http.NewServeMux()
	empHandlers := NewEmployeeHandlers(empService)
	deptHandlers := NewDepartmentHandlers(deptService)

	// Employee routes
	router.HandleFunc("POST /employees", empHandlers.CreateEmployee)
	router.HandleFunc("GET /employees/{id}", empHandlers.GetEmployee)
	router.HandleFunc("PUT /employees/{id}", empHandlers.UpdateEmployee)
	router.HandleFunc("DELETE /employees/{id}", empHandlers.DeleteEmployee)
	router.HandleFunc("GET /companies/{companyId}/employees", empHandlers.GetCompanyEmployees)
	router.HandleFunc("GET /companies/{companyId}/departments/{departmentId}/employees", empHandlers.GetDepartmentEmployees)

	// Department routes
	router.HandleFunc("POST /departments", deptHandlers.GetOrCreateDepartment)
	router.HandleFunc("GET /departments/{id}", deptHandlers.GetDepartment)

	return router
}
