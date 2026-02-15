package models

import "time"

type AssignmentStatus string

const (
	StatusAssigned  AssignmentStatus = "Assigned"
	StatusReleased  AssignmentStatus = "Released"
	StatusRejected  AssignmentStatus = "Rejected"
	StatusRequested AssignmentStatus = "Requested"
)

// Shift represents a time slot created by a supervisor
type Shift struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Creator Employee `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// ShiftAssignment represents a worker assigned to a shift
type ShiftAssignment struct {
	ID         uint             `gorm:"primaryKey" json:"id"`
	ShiftID    uint             `json:"shift_id"`
	EmployeeID uint             `json:"employee_id"`
	AssigneeID uint             `json:"assigned_by"`
	AssignedAt time.Time        `json:"assigned_at"`
	Status     AssignmentStatus `json:"status"`

	// Relations
	Shift    Shift    `gorm:"foreignKey:ShiftID" json:"shift,omitempty"`
	Employee Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`
	Assignee Employee `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
}
