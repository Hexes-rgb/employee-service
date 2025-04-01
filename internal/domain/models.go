package domain

type Department struct {
	ID        int    `json:"id"`
	CompanyID int    `json:"companyId"`
	Name      string `json:"name"`
	Phone     string `json:"phone,omitempty"`
}

type Employee struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	Surname        string      `json:"surname"`
	Phone          string      `json:"phone,omitempty"`
	CompanyID      int         `json:"companyId"`
	DepartmentID   *int        `json:"departmentId,omitempty"`
	PassportType   string      `json:"passportType,omitempty"`
	PassportNumber string      `json:"passportNumber,omitempty"`
	Department     *Department `json:"department,omitempty"`
}
