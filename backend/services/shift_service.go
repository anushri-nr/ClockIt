package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"errors"
	"log"
	"time"
)

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
