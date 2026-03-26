package controllers

import (
	"clockit/backend/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateShift creates a new shift time slot
func CreateShift(c *gin.Context) {
	type Req struct {
		StartTime string `json:"start_time" binding:"required"`
		EndTime   string `json:"end_time" binding:"required"`
<<<<<<< HEAD
=======
		CreatedBy uint   `json:"created_by" binding:"required"`
>>>>>>> main
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateShift: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

<<<<<<< HEAD
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

=======
>>>>>>> main
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

<<<<<<< HEAD
	log.Printf("Creating shift: start=%s, end=%s, creator=%d", req.StartTime, req.EndTime, createdBy)

	shift, err := services.CreateShift(startTime, endTime, createdBy)
=======
	log.Printf("Creating shift: start=%s, end=%s, creator=%d", req.StartTime, req.EndTime, req.CreatedBy)

	shift, err := services.CreateShift(startTime, endTime, req.CreatedBy)
>>>>>>> main
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
<<<<<<< HEAD
		// AssignedBy is removed! The frontend no longer needs to send it.
=======
		AssignedBy uint `json:"assigned_by" binding:"required"` // Temporary field to track who made the assignment (supervisor ID) for auditing purposes. Replace with authenticated user context.
>>>>>>> main
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("AssignWorkerToShift: binding error: %v", err)
<<<<<<< HEAD
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
=======
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	log.Printf("Assigning worker: shift=%d, employee=%d, assignedBy=%d", req.ShiftID, req.EmployeeID, req.AssignedBy)

	assignment, err := services.AssignWorkerToShift(req.ShiftID, req.EmployeeID, req.AssignedBy)
>>>>>>> main
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
