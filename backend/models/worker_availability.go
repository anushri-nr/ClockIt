package models

import "time"

type WorkerAvailability struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	WorkerID  uint      `gorm:"constraint:OnDelete:CASCADE" json:"worker_id"`
	DayOfWeek int       `json:"day_of_week"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Worker Employee `gorm:"foreignKey:WorkerID" json:"worker,omitempty"`
}
