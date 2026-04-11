package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// ErrAssignmentNotAssigned is returned when attempting to release an assignment
// that is not currently in the Assigned state.
var ErrAssignmentNotAssigned = errors.New("assignment must be in Assigned status to be released")

// CreateShift creates a new shift time slot
func CreateShift(startTime, endTime time.Time, createdBy uint) (*models.Shift, error) {
	log.Printf("CreateShift: creator=%d, start=%s, end=%s", createdBy, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))

	// Verify creator exists and is a supervisor
	var supervisor models.Employee
	if err := database.DB.Where("id = ? AND role = ?", createdBy, models.RoleSupervisor).First(&supervisor).Error; err != nil {
		log.Printf("CreateShift: supervisor not found: %v", err)
		return nil, errors.New("supervisor not found")
	}

	// Validate time order
	if !endTime.After(startTime) {
		return nil, errors.New("end_time must be after start_time")
	}

	shift := models.Shift{
		StartTime: startTime,
		EndTime:   endTime,
		CreatedBy: createdBy,
	}

	if err := database.DB.Create(&shift).Error; err != nil {
		log.Printf("CreateShift: DB insert failed: %v", err)
		return nil, err
	}

	// Load creator relation
	if err := database.DB.Preload("Creator").First(&shift, shift.ID).Error; err != nil {
		log.Printf("CreateShift: failed to load creator relation: %v", err)
		// Assignment created successfully, but relations failed to load
		// Return the shift without relations rather than failing the entire operation
	}

	log.Printf("CreateShift: shift created, id=%d", shift.ID)
	return &shift, nil
}

// AssignWorkerToShift assigns a worker to an existing shift
func AssignWorkerToShift(shiftID, employeeID, assignedBy uint) (*models.ShiftAssignment, error) {
	log.Printf("AssignWorkerToShift: shift=%d, employee=%d, assignedBy=%d", shiftID, employeeID, assignedBy)

	// Verify shift exists
	var shift models.Shift
	if err := database.DB.First(&shift, shiftID).Error; err != nil {
		log.Printf("AssignWorkerToShift: shift not found: %v", err)
		return nil, errors.New("shift not found")
	}

	// Verify employee exists and is a worker
	var worker models.Employee
	if err := database.DB.Where("id = ? AND role = ?", employeeID, models.RoleWorker).First(&worker).Error; err != nil {
		log.Printf("AssignWorkerToShift: worker not found: %v", err)
		return nil, errors.New("worker not found")
	}

	// Verify assigner is a supervisor
	var supervisor models.Employee
	if err := database.DB.Where("id = ? AND role = ?", assignedBy, models.RoleSupervisor).First(&supervisor).Error; err != nil {
		log.Printf("AssignWorkerToShift: supervisor not found: %v", err)
		return nil, errors.New("supervisor not found or insufficient permissions")
	}

	assignment := models.ShiftAssignment{
		ShiftID:    shiftID,
		EmployeeID: employeeID,
		AssigneeID: assignedBy,
		Status:     models.StatusAssigned,
		AssignedAt: time.Now(),
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		log.Printf("AssignWorkerToShift: DB insert failed: %v", err)
		return nil, err
	}

	// Load relations
	if err := database.DB.Preload("Shift").Preload("Employee").Preload("Assignee").First(&assignment, assignment.ID).Error; err != nil {
		log.Printf("AssignWorkerToShift: failed to load relations: %v", err)
		// Assignment created successfully, but relations failed to load
		// Return the assignment without relations rather than failing the entire operation
	}

	log.Printf("AssignWorkerToShift: assignment created, id=%d", assignment.ID)
	return &assignment, nil
}

// ReleaseShiftForWorker marks an existing assignment for a worker as Released
func ReleaseShiftForWorker(shiftID, employeeID uint) (*models.ShiftAssignment, error) {
	log.Printf("ReleaseShiftForWorker: shift=%d, employee=%d", shiftID, employeeID)

	var assignment models.ShiftAssignment
	if err := database.DB.Where("shift_id = ? AND employee_id = ?", shiftID, employeeID).First(&assignment).Error; err != nil {
		log.Printf("ReleaseShiftForWorker: assignment lookup error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Only allow releasing assignments that are currently Assigned
	if assignment.Status != models.StatusAssigned {
		return nil, ErrAssignmentNotAssigned
	}

	assignment.Status = models.StatusReleased

	if err := database.DB.Save(&assignment).Error; err != nil {
		log.Printf("ReleaseShiftForWorker: failed to update assignment: %v", err)
		return nil, err
	}

	// Attempt to load relations. If it fails, we still return the assignment
	if err := database.DB.Preload("Shift").Preload("Employee").Preload("Assignee").First(&assignment, assignment.ID).Error; err != nil {
		log.Printf("ReleaseShiftForWorker: failed to preload relations: %v", err)
	}

	log.Printf("ReleaseShiftForWorker: assignment updated, id=%d", assignment.ID)
	return &assignment, nil
}

// RejectShiftRequest lets a supervisor reject a requested shift in their company.
// Status transition: Requested -> Released
func (s *SupervisorService) RejectShiftRequest(supervisorID, shiftID uint) (*models.ShiftAssignment, error) {
	var out models.ShiftAssignment

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// validate supervisor
		var sup models.Employee
		if err := tx.Where("id = ? AND role = ?", supervisorID, models.RoleSupervisor).First(&sup).Error; err != nil {
			return err
		}

		// find requested assignment for shift in supervisor's company
		var a models.ShiftAssignment
		err := tx.
			Table("shift_assignments as sa").
			Select("sa.*").
			Joins("JOIN shifts s ON s.id = sa.shift_id").
			Joins("JOIN employees creator ON creator.id = s.created_by").
			Where("sa.shift_id = ?", shiftID).
			Where("sa.status = ?", models.StatusRequested).
			Where("creator.company_id = ?", sup.CompanyID).
			First(&a).Error
		if err != nil {
			return err
		}

		a.Status = models.StatusReleased
		a.AssigneeID = supervisorID
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

// // RequestReleasedShift lets a worker request a released shift.
// // It updates shift.status -> requested and creates/updates shift_assignments.
// func RequestReleasedShift(workerID, shiftID uint) (*models.ShiftAssignment, error) {
// 	var out models.ShiftAssignment

// 	err := database.DB.Transaction(func(tx *gorm.DB) error {
// 		var worker models.Employee
// 		if err := tx.Where("id = ? AND role = ?", workerID, models.RoleWorker).First(&worker).Error; err != nil {
// 			return err
// 		}

// 		// find released assignment row for this shift
// 		var a models.ShiftAssignment
// 		if err := tx.
// 			Where("shift_id = ? AND status = ?", shiftID, models.StatusReleased).
// 			First(&a).Error; err != nil {
// 			return err
// 		}

// 		// optional: ensure shift belongs to worker company
// 		var shift models.Shift
// 		if err := tx.First(&shift, shiftID).Error; err != nil {
// 			return err
// 		}
// 		var creator models.Employee
// 		if err := tx.First(&creator, shift.CreatedBy).Error; err != nil {
// 			return err
// 		}
// 		if creator.CompanyID != worker.CompanyID {
// 			return fmt.Errorf("shift is not from worker company")
// 		}

// 		// update assignment -> requested by this worker
// 		a.EmployeeID = workerID
// 		a.Status = models.StatusRequested
// 		a.AssignedAt = time.Now().UTC()
// 		if err := tx.Save(&a).Error; err != nil {
// 			return err
// 		}

// 		out = a
// 		return nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &out, nil
// }

func RequestReleasedShift(workerID, shiftID uint) (*models.ShiftAssignment, error) {
	var assignment models.ShiftAssignment

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// worker must exist
		var worker models.Employee
		if err := tx.Where("id = ? AND role = ?", workerID, models.RoleWorker).First(&worker).Error; err != nil {
			return err
		}

		// find the existing released assignment row
		if err := tx.
			Where("shift_id = ? AND employee_id = ? AND status = ?", shiftID, workerID, models.StatusReleased).
			First(&assignment).Error; err != nil {
			return err
		}

		// update only the status
		if err := tx.Model(&assignment).
			Update("status", models.StatusRequested).Error; err != nil {
			return err
		}

		assignment.Status = models.StatusRequested
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &assignment, nil
}
