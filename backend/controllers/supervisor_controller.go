package controllers

import (
	"clockit/backend/services"
	"log"
	"net/http"
	"time"

	"clockit/backend/models"

	"github.com/gin-gonic/gin"
)

func SupervisorRegister(c *gin.Context) {
	type Req struct {
		Name      string  `json:"name" binding:"required"`
		Email     string  `json:"email" binding:"required,email"`
		Password  string  `json:"password" binding:"required,min=8,max=20"`
		Address   string  `json:"address"`
		PhoneNo   string  `json:"phone_no"`
		CompanyID uint    `json:"company_id"`
		Wage      float64 `json:"wage"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("SupervisorRegister: invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
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
		"id":        employee.ID,
		"name":      employee.Name,
		"email":     employee.Email,
		"role":      employee.Role,
		"companyID": employee.CompanyID,
	})
}

// Request body for worker availability
type AvailabilityRequest struct {
	Date string `json:"date" binding:"required"`
}

// Returns all workers available on a specific date
func GetWorkerAvailability(c *gin.Context) {
	var req AvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("GetWorkerAvailability: invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request, date required"})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		log.Printf("GetWorkerAvailability: invalid date format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	workers, err := services.GetWorkersAvailableForDate(date)
	if err != nil {
		log.Printf("GetWorkerAvailability: failed to fetch workers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch availability"})
		return
	}

	log.Printf("Found %d workers available", len(workers))
	c.JSON(http.StatusOK, workers)
}
