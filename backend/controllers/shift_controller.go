package controllers

import (
	"clockit/backend/services"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// CreateShift creates a new shift time slot
func CreateShift(c *gin.Context) {
	type Req struct {
		StartTime string `json:"start_time" binding:"required"`
		EndTime   string `json:"end_time" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateShift: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	// Extract the ID directly from the validated JWT Token
	authIDVal, exists := c.Get(services.ContextEmployeeID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}
	authID, ok := authIDVal.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}
	createdBy := authID

	// Parse timestamps (ISO 8601 format)
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		log.Printf("CreateShift: invalid start_time format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format, use ISO 8601 (e.g., 2026-02-15T09:00:00Z)"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		log.Printf("CreateShift: invalid end_time format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format, use ISO 8601 (e.g., 2026-02-15T17:00:00Z)"})
		return
	}

	log.Printf("Creating shift: start=%s, end=%s, creator=%d", req.StartTime, req.EndTime, createdBy)

	shift, err := services.CreateShift(startTime, endTime, createdBy)
	if err != nil {
		log.Printf("Failed to create shift: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Shift created successfully, id=%d", shift.ID)
	c.JSON(http.StatusCreated, gin.H{
		"id":         shift.ID,
		"start_time": shift.StartTime.Format(time.RFC3339),
		"end_time":   shift.EndTime.Format(time.RFC3339),
		"created_by": shift.CreatedBy,
		"created_at": shift.CreatedAt.Format(time.RFC3339),
	})
}

// AssignWorkerToShift assigns a worker to an existing shift
func AssignWorkerToShift(c *gin.Context) {
	type Req struct {
		ShiftID    uint `json:"shift_id" binding:"required"`
		EmployeeID uint `json:"employee_id" binding:"required"`
		// AssignedBy is removed! The frontend no longer needs to send it.
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("AssignWorkerToShift: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload provided"})
		return
	}

	// Extract the ID directly from the validated JWT Token
	authIDVal, exists := c.Get(services.ContextEmployeeID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}
	assignedBy := authIDVal.(uint) // Safely cast the token ID to a uint

	log.Printf("Assigning worker: shift=%d, employee=%d, assignedBy=%d", req.ShiftID, req.EmployeeID, assignedBy)

	assignment, err := services.AssignWorkerToShift(req.ShiftID, req.EmployeeID, assignedBy)
	if err != nil {
		log.Printf("Failed to assign worker: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Worker assigned successfully, assignment_id=%d", assignment.ID)
	c.JSON(http.StatusCreated, gin.H{
		"id":          assignment.ID,
		"shift_id":    assignment.ShiftID,
		"employee_id": assignment.EmployeeID,
		"assigned_by": assignment.AssigneeID,
		"status":      assignment.Status,
		"assigned_at": assignment.AssignedAt.Format(time.RFC3339),
	})
}

// ReleaseShiftForWorker allows a worker to release a shift they were assigned to
func ReleaseShiftForWorker(c *gin.Context) {
	type Req struct {
		ShiftID uint `json:"shift_id" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ReleaseShiftForWorker: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	// get authenticated employee id from JWT
	empID, err := services.GetAuthenticatedEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	log.Printf("Worker %d releasing shift %d", empID, req.ShiftID)

	assignment, err := services.ReleaseShiftForWorker(req.ShiftID, empID)
	if err != nil {
		log.Printf("ReleaseShiftForWorker: service error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "shift or assignment not found"})
			return
		}
		if errors.Is(err, services.ErrAssignmentNotAssigned) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "assignment must be in Assigned status to be released"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to release assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          assignment.ID,
		"shift_id":    assignment.ShiftID,
		"employee_id": assignment.EmployeeID,
		"status":      assignment.Status,
		"assigned_at": assignment.AssignedAt.Format(time.RFC3339),
	})
}

// RejectShiftRequest rejects a requested shift and re-releases it.
func RejectShiftRequest(c *gin.Context) {
	authID, err := services.GetAuthenticatedEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	shiftID64, err := strconv.ParseUint(c.Param("shift_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift_id"})
		return
	}

	svc := services.SupervisorService{}
	assignment, err := svc.RejectShiftRequest(authID, uint(shiftID64))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "requested shift not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject shift request"})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

// RequestShift lets an authenticated worker request a released shift.
func RequestShift(c *gin.Context) {
	workerID, err := services.GetAuthenticatedEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	shiftID64, err := strconv.ParseUint(c.Param("shift_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift_id"})
		return
	}

	assignment, err := services.RequestReleasedShift(workerID, uint(shiftID64))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "worker or released shift not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

// DeleteShift completely removes a shift
func DeleteShift(c *gin.Context) {
	// Extract the ID directly from the validated JWT Token
	authIDVal, exists := c.Get(services.ContextEmployeeID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}
	supervisorID := authIDVal.(uint)

	// Extract the shift_id from the URL
	shiftID64, err := strconv.ParseUint(c.Param("shift_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift_id"})
		return
	}

	// Call the service
	if err := services.DeleteShift(uint(shiftID64), supervisorID); err != nil {
		log.Printf("Failed to delete shift: %v", err)
		if err.Error() == "unauthorized to delete shifts outside your company" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete shift"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shift deleted successfully"})
}

// UnassignWorker removes a worker from a shift
func UnassignWorker(c *gin.Context) {
	// Extract supervisor ID from JWT
	authIDVal, exists := c.Get(services.ContextEmployeeID)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}
	supervisorID := authIDVal.(uint)

	// Extract shift_id from URL
	shiftID64, err := strconv.ParseUint(c.Param("shift_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shift_id"})
		return
	}

	if err := services.UnassignWorkerFromShift(uint(shiftID64), supervisorID); err != nil {
		log.Printf("Failed to unassign worker: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
 			c.JSON(http.StatusNotFound, gin.H{"error": "shift not found"})
 			return
 		}
 		if err.Error() == "unauthorized to unassign shifts outside your company" {
 			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
 			return
 		}
 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unassign worker"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "worker unassigned successfully"})
}