package models

type EmployeeRole string

const (
	RoleSupervisor EmployeeRole = "supervisor"
	RoleWorker     EmployeeRole = "worker"
)

type Employee struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	CompanyID uint         `json:"company_id" binding:"required"`
	Name      string       `json:"name"`
	Address   string       `json:"address"`
	Email     string       `gorm:"unique" json:"email"`
	PhoneNo   string       `json:"phone_no"`
	Role      EmployeeRole `json:"role"`
	Wage      float64      `json:"wage" binding:"gte=0" gorm:"type:decimal(10,2)"`
	Password  string       `json:"-"`

	// Relations
	Company Company `gorm:"foreignKey:CompanyID;constraint:OnDelete:RESTRICT" json:"company,omitempty"`
}
