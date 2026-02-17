package controllers

import (
	"clockit/backend/models"
	"clockit/backend/repository"
	"clockit/backend/services"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WorkerRegister(c *gin.Context) {
	type Req struct {
		Name      string  `json:"name" binding:"required"`
		Email     string  `json:"email" binding:"required,email"`
		Password  string  `json:"password" binding:"required,min=8,max=20"`
		Address   string  `json:"address"`
		PhoneNo   string  `json:"phone_no"`
		CompanyID uint    `json:"company_id" binding:"required"`
		Wage      float64 `json:"wage" binding:"gte=0"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("WorkerRegister: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	log.Printf("Registering worker %s (%s)", req.Name, req.Email)

	employee, err := services.RegisterEmployee(
		req.Name,
		req.Email,
		req.Password,
		req.Address,
		req.PhoneNo,
		req.CompanyID,
		req.Wage,
		models.RoleWorker,
	)
	if err != nil {
		log.Printf("Failed to register worker: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Worker registered successfully, id=%d", employee.ID)
	c.JSON(http.StatusCreated, gin.H{
		"id":        employee.ID,
		"name":      employee.Name,
		"email":     employee.Email,
		"role":      employee.Role,
		"companyID": employee.CompanyID,
	})
}

// Query params for fetching assigned shifts
// GetShiftsForWorker returns all shifts assigned to a worker
func GetShiftsForWorker(c *gin.Context) {
	// parse employee_id from path
	empParam := c.Param("employee_id")
	if empParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id path parameter is required"})
		return
	}

	// convert empParam to uint
	empID64, err := strconv.ParseUint(empParam, 10, 64)
	if err != nil {
		log.Printf("GetShiftsForWorker: invalid employee_id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}

	svc := services.WorkerService{
		Repo: &repository.WorkerRepository{},
	}

	resp, err := svc.GetAssignedShifts(uint(empID64))
	if err != nil {
		log.Printf("GetShiftsForWorker: service error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "worker not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch assigned shifts"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
