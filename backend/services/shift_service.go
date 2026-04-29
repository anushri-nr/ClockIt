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

	// 1. Force an UPDATE if the record already exists (e.g., from a Request)
	result := database.DB.Table("shift_assignments").
		Where("shift_id = ? AND employee_id = ?", shiftID, employeeID).
		Updates(map[string]interface{}{
			"assignee_id": assignedBy,
			"status":      models.StatusAssigned,
			"assigned_at": time.Now(),
		})

	if result.Error != nil {
		log.Printf("AssignWorkerToShift: DB update error: %v", result.Error)
		return nil, result.Error
	}

	// 2. If RowsAffected is 0, the record didn't exist, so we MUST insert it.
	if result.RowsAffected == 0 {
		newAssignment := models.ShiftAssignment{
			ShiftID:    shiftID,
			EmployeeID: employeeID,
			AssigneeID: assignedBy,
			Status:     models.StatusAssigned,
			AssignedAt: time.Now(),
		}
		if err := database.DB.Create(&newAssignment).Error; err != nil {
			log.Printf("AssignWorkerToShift: DB insert failed: %v", err)
			return nil, err
		}
	}

	// 3. Fetch the final record to return and check for errors
	var assignment models.ShiftAssignment
	if err := database.DB.Preload("Shift").Preload("Employee").Preload("Assignee").
		Where("shift_id = ? AND employee_id = ?", shiftID, employeeID).First(&assignment).Error; err != nil {
		log.Printf("AssignWorkerToShift: failed to fetch final assignment: %v", err)
		return nil, err
	}

	log.Printf("AssignWorkerToShift: success for shift=%d", shiftID)
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
	var assignment models.ShiftAssignment

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var sup models.Employee
		if err := tx.Where("id = ? AND role = ?", supervisorID, models.RoleSupervisor).First(&sup).Error; err != nil {
			return err
		}

		if err := tx.
			Table("shift_assignments as sa").
			Select("sa.*").
			Joins("JOIN shifts s ON s.id = sa.shift_id").
			Joins("JOIN employees creator ON creator.id = s.created_by").
			Where("sa.shift_id = ?", shiftID).
			Where("sa.status = ?", models.StatusRequested).
			Where("creator.company_id = ?", sup.CompanyID).
			First(&assignment).Error; err != nil {
			return err
		}

		if err := tx.Model(&assignment).Updates(map[string]interface{}{
			"status":      models.StatusReleased,
			"assignee_id": supervisorID,
			"assigned_at": time.Now().UTC(),
		}).Error; err != nil {
			return err
		}

		assignment.Status = models.StatusReleased
		assignment.AssigneeID = supervisorID
		assignment.AssignedAt = time.Now().UTC()
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &assignment, nil
}


func RequestReleasedShift(workerID, shiftID uint) (*models.ShiftAssignment, error) {
	var assignment models.ShiftAssignment

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// validate worker
		var worker models.Employee
		if err := tx.
			Where("id = ? AND role = ?", workerID, models.RoleWorker).
			First(&worker).Error; err != nil {
			return err
		}

		// find one released assignment for this shift in worker's company
		if err := tx.
			Table("shift_assignments AS sa").
			Select("sa.*").
			Joins("JOIN shifts s ON s.id = sa.shift_id").
			Joins("JOIN employees creator ON creator.id = s.created_by").
			Where("sa.shift_id = ?", shiftID).
			Where("sa.status = ?", models.StatusReleased).
			Where("creator.company_id = ?", worker.CompanyID).
			Order("sa.assigned_at DESC").
			First(&assignment).Error; err != nil {
			return err
		}

		// update same row (no new row)
		now := time.Now().UTC()
		if err := tx.Model(&assignment).Updates(map[string]interface{}{
			"employee_id": workerID,
			"status":      models.StatusRequested,
			"assigned_at": now,
		}).Error; err != nil {
			return err
		}

		assignment.EmployeeID = workerID
		assignment.Status = models.StatusRequested
		assignment.AssignedAt = now
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

// DeleteShift completely removes a shift and any associated assignments
func DeleteShift(shiftID uint, supervisorID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Verify supervisor exists
		var sup models.Employee
		if err := tx.Where("id = ? AND role = ?", supervisorID, models.RoleSupervisor).First(&sup).Error; err != nil {
			return err // Bubbles up gorm.ErrRecordNotFound
		}

		// 2. Verify shift exists
		var shift models.Shift
		if err := tx.Preload("Assignments").First(&shift, shiftID).Error; err != nil {
			return err // Bubbles up gorm.ErrRecordNotFound
		}

		// Guard to prevent deleting assigned shifts
		for _, a := range shift.Assignments {
			if a.Status == models.StatusAssigned || a.Status == models.StatusRequested {
				return errors.New("cannot delete a shift that has active assignments")
			}
		}

		return nil
	})
}

// UnassignWorkerFromShift removes the assigned worker from a shift
func UnassignWorkerFromShift(shiftID, supervisorID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Verify supervisor
		var sup models.Employee
		if err := tx.Where("id = ? AND role = ?", supervisorID, models.RoleSupervisor).First(&sup).Error; err != nil {
			return errors.New("supervisor not found")
		}

		// 2. Verify shift exists and belongs to the supervisor's company
		var shift models.Shift
		if err := tx.First(&shift, shiftID).Error; err != nil {
			return err 
		}
		if shift.Creator.CompanyID != sup.CompanyID {
			return errors.New("unauthorized to unassign shifts outside your company")
		}

		// 3. Delete the active assignment row so it returns to Unassigned
		result := tx.Unscoped().Where("shift_id = ? AND status = ?", shiftID, models.StatusAssigned).Delete(&models.ShiftAssignment{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("no active assignment found for this shift")
		}

		return nil
	})
}
