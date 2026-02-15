package models

import "time"

type Company struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	// Back-reference to employees in this company (not stored in DB, used for eager loading)
	Employees []Employee `gorm:"foreignKey:CompanyID" json:"employees,omitempty"`
}
