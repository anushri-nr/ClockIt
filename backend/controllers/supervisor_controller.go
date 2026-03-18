package controllers

import (
	"clockit/backend/repository"
	"clockit/backend/services"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"clockit/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SupervisorRegister(c *gin.Context) {
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
		log.Printf("SupervisorRegister: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	log.Printf("Registering supervisor %s (%s)", req.Name, req.Email)

	employee, err := services.RegisterEmployee(
		req.Name,
		req.Email,
		req.Password,
		req.Address,
		req.PhoneNo,
		req.CompanyID,
		req.Wage,
		models.RoleSupervisor,
	)
	if err != nil {
		log.Printf("Failed to register supervisor: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Supervisor registered successfully, id=%d", employee.ID)
	c.JSON(http.StatusCreated, gin.H{
		"id":         employee.ID,
		"name":       employee.Name,
		"email":      employee.Email,
		"role":       employee.Role,
		"company_id": employee.CompanyID,
	})
}

// Request body for worker availability
type AvailabilityQuery struct {
	Date      string `form:"date" binding:"required"`
	CompanyID uint   `form:"company_id" binding:"required"`
}

// Returns all workers available on a specific date
func GetWorkerAvailability(c *gin.Context) {
	var query AvailabilityQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		log.Printf("GetWorkerAvailability: invalid query: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date and company_id query parameters are required",
		})
		return
	}

	date, err := time.Parse("2006-01-02", query.Date)
	if err != nil {
		log.Printf("GetWorkerAvailability: invalid date format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	service := services.WorkerService{
		Repo: &repository.WorkerRepository{},
	}

	workers, err := service.GetWorkersAvailable(date, query.CompanyID)
	if err != nil {
		log.Printf("GetWorkerAvailability: failed to fetch workers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch availability"})
		return
	}

	log.Printf("Found %d workers available", len(workers))
	c.JSON(http.StatusOK, workers)
}

// GetShiftsByCompany returns shifts for a supervisor's company, optionally filtered by status
func GetShiftsByCompany(c *gin.Context) {
	empParam := c.Param("employee_id")
	if empParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id path parameter is required"})
		return
	}

	statusParam := c.Query("status")
	
	// Validate status against the supported set
	if statusParam != "" {
		validStatuses := map[string]bool{
			string(models.StatusAssigned):   true,
			string(models.StatusReleased):   true,
			string(models.StatusRejected):   true,
			string(models.StatusRequested):  true,
			string(models.StatusUnassigned): true,
		}
		if !validStatuses[statusParam] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status filter provided"})
			return
		}
	}

	empID64, err := strconv.ParseUint(empParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}

	svc := services.SupervisorService{}
	
	shifts, err := svc.GetShiftsByCompany(uint(empID64), statusParam)
	if err != nil {
		log.Printf("GetShiftsByCompany: service error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "supervisor not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch shifts"})
		return
	}

	c.JSON(http.StatusOK, shifts)
}