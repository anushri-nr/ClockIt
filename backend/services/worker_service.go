package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"errors"
	"fmt"
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

// GetReleasedShiftsByEmployeeCompany returns released shifts from the employee's company.
func (s *WorkerService) GetReleasedShiftsByEmployeeCompany(employeeID uint) ([]models.Shift, error) {
    var emp models.Employee
    if err := database.DB.
        Select("id", "company_id").
        Where("id = ? AND role = ?", employeeID, models.RoleWorker).
        First(&emp).Error; err != nil {
        return []models.Shift{}, err
    }

    var shifts []models.Shift
    err := database.DB.
        Model(&models.Shift{}).
        Select("DISTINCT shifts.*").
        Joins("JOIN shift_assignments sa ON sa.shift_id = shifts.id").
        Joins("JOIN employees creator ON creator.id = shifts.created_by").
        Where("creator.company_id = ?", emp.CompanyID).
        Where("sa.status = ?", models.StatusReleased).
        Order("shifts.start_time ASC").
        Find(&shifts).Error
    if err != nil {
        return []models.Shift{}, err
    }

    return shifts, nil
}

// RequestReleasedShift lets a worker request a released shift.
// It updates shift.status -> requested and creates/updates shift_assignments.
func RequestReleasedShift(workerID, shiftID uint) (*models.ShiftAssignment, error) {
    var out models.ShiftAssignment

    err := database.DB.Transaction(func(tx *gorm.DB) error {
        var worker models.Employee
        if err := tx.Where("id = ? AND role = ?", workerID, models.RoleWorker).First(&worker).Error; err != nil {
            return err
        }

        // find released assignment row for this shift
        var a models.ShiftAssignment
        if err := tx.
            Where("shift_id = ? AND status = ?", shiftID, models.StatusReleased).
            First(&a).Error; err != nil {
            return err
        }

        // optional: ensure shift belongs to worker company
        var shift models.Shift
        if err := tx.First(&shift, shiftID).Error; err != nil {
            return err
        }
        var creator models.Employee
        if err := tx.First(&creator, shift.CreatedBy).Error; err != nil {
            return err
        }
        if creator.CompanyID != worker.CompanyID {
            return fmt.Errorf("shift is not from worker company")
        }

        // update assignment -> requested by this worker
        a.EmployeeID = workerID
        a.Status = models.StatusRequested
        a.AssignedAt = time.Now().UTC()
        if err := tx.Save(&a).Error; err != nil {
            return err
        }

        out = a
        return nil
    })

	if err != nil {
        return nil, err
    }
    return &out, nil
}