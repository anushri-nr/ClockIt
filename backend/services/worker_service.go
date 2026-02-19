package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// AssignedShiftResponse is the public response shape for assigned shifts
type AssignedShiftResponse struct {
	ID         uint   `json:"id"`
	ShiftID    uint   `json:"shift_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	AssignedBy uint   `json:"assigned_by"`
	AssignedAt string `json:"assigned_at"`
	Status     string `json:"status"`
}

// GetAssignedShifts returns assigned shifts for a worker
func (s *WorkerService) GetAssignedShifts(employeeID uint) ([]AssignedShiftResponse, error) {
	// Validate employee exists and is a worker
	var emp models.Employee
	if err := database.DB.Where("id = ? AND role = ?", employeeID, models.RoleWorker).First(&emp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("GetAssignedShifts: employee not found or not a worker: %v", err)
			return []AssignedShiftResponse{}, gorm.ErrRecordNotFound
		}
		log.Printf("GetAssignedShifts: DB error while validating employee: %v", err)
		return []AssignedShiftResponse{}, err
	}

	assignments, err := s.Repo.FindAssignmentsForWorker(employeeID)
	if err != nil {
		log.Printf("GetAssignedShifts: failed to fetch assignments from repository: %v", err)
		return []AssignedShiftResponse{}, err
	}

	resp := make([]AssignedShiftResponse, 0)
	for _, a := range assignments {
		r := AssignedShiftResponse{
			ID:         a.ID,
			ShiftID:    a.ShiftID,
			AssignedBy: a.AssigneeID,
			Status:     string(a.Status),
		}
		if !a.AssignedAt.IsZero() {
			r.AssignedAt = a.AssignedAt.Format(time.RFC3339)
		}
		if a.Shift.ID != 0 {
			r.StartTime = a.Shift.StartTime.Format(time.RFC3339)
			r.EndTime = a.Shift.EndTime.Format(time.RFC3339)
		}
		resp = append(resp, r)
	}

	log.Printf("GetAssignedShifts: returning %d assignments for employee=%d", len(resp), employeeID)
	return resp, nil
}
