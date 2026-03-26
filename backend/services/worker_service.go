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
	ID            uint    `json:"id"`
	ShiftID       uint    `json:"shift_id"`
	StartTime     string  `json:"start_time"`
	EndTime       string  `json:"end_time"`
	AssignedBy    uint    `json:"assigned_by"`
	AssignedAt    string  `json:"assigned_at"`
	Status        string  `json:"status"`
	DurationHours float64 `json:"duration_hours"` // Added for Wage Calc
	Earnings      float64 `json:"earnings"`       // Added for Wage Calc
}

// GetAssignedShifts returns assigned shifts for a worker, including calculated wages
func (s *WorkerService) GetAssignedShifts(employeeID uint) ([]AssignedShiftResponse, error) {
	// Validate employee exists and is a worker, and load their wage
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
			
			// Issue #38: Wage Calculation
			duration := a.Shift.EndTime.Sub(a.Shift.StartTime).Hours()
			r.DurationHours = duration
			r.Earnings = duration * emp.Wage
		}
		
		resp = append(resp, r)
	}

	log.Printf("GetAssignedShifts: returning %d assignments for employee=%d", len(resp), employeeID)
	return resp, nil
}

// CreateWorkerAvailability allows a worker to set their availability for a specific day of the week
func CreateWorkerAvailability(workerID uint, dayOfWeek int, startTime string, endTime string) (*models.WorkerAvailability, error) {
	// Validate worker exists and is a worker
	var emp models.Employee
	if err := database.DB.Where("id = ? AND role = ?", workerID, models.RoleWorker).First(&emp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("CreateWorkerAvailability: employee not found or not a worker: %v", err)
			return nil, gorm.ErrRecordNotFound
		}
		log.Printf("CreateWorkerAvailability: DB error while validating employee: %v", err)
		return nil, err
	}

	availability := &models.WorkerAvailability{
		WorkerID:  workerID,
		DayOfWeek: dayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := database.DB.Create(availability).Error; err != nil {
		log.Printf("CreateWorkerAvailability: failed to create availability in DB: %v", err)
		return nil, err
	}

	log.Printf("CreateWorkerAvailability: availability created successfully for worker=%d on day=%d", workerID, dayOfWeek)
	return availability, nil
}
