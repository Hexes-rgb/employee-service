package domain

type Department struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"companyId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type Employee struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	Surname        string      `json:"surname"`
	Phone          string      `json:"phone"`
	CompanyID      int         `json:"companyId"`
	DepartmentID   *int        `json:"departmentId"`
	PassportType   string      `json:"passportType"`
	PassportNumber string      `json:"passportNumber"`
	Department     *Department `json:"department"`
}
