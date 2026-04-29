package controllers

import (
	"clockit/backend/models"
	"clockit/backend/repository"
	"clockit/backend/services"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

	token, err := services.GenerateToken(employee.ID, string(employee.Role))
	if err != nil {
		log.Printf("WorkerRegister: token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	log.Printf("Worker registered successfully, id=%d", employee.ID)
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"employee": gin.H{
			"id":         employee.ID,
			"name":       employee.Name,
			"email":      employee.Email,
			"role":       employee.Role,
			"company_id": employee.CompanyID,
		},
	})
}

// WorkerLogin authenticates a worker and returns a JWT token
func WorkerLogin(c *gin.Context) {
	type Req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("WorkerLogin: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	// find worker by email
	var worker models.Employee
	if err := services.FindEmployeeByEmailAndRole(req.Email, models.RoleWorker, &worker); err != nil {
		log.Printf("WorkerLogin: lookup failed: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		}
		return
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(worker.Password), []byte(req.Password)); err != nil {
		log.Printf("WorkerLogin: password mismatch for email=%s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// generate JWT
	token, err := services.GenerateToken(worker.ID, string(worker.Role))
	if err != nil {
		log.Printf("WorkerLogin: token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"employee": gin.H{
			"id":         worker.ID,
			"name":       worker.Name,
			"email":      worker.Email,
			"role":       worker.Role,
			"company_id": worker.CompanyID,
		},
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

	var startTime time.Time
	var endTime time.Time
	var filterByWindow bool

	startParam := c.Query("start_time")
	endParam := c.Query("end_time")
	if startParam != "" || endParam != "" {
		if startParam == "" || endParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "both start_time and end_time must be provided"})
			return
		}
		var errP error
		startTime, errP = time.Parse(time.RFC3339, startParam)
		if errP != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format, use RFC3339"})
			return
		}
		endTime, errP = time.Parse(time.RFC3339, endParam)
		if errP != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format, use RFC3339"})
			return
		}
		if endTime.Before(startTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "end_time must be equal or after start_time"})
			return
		}
		filterByWindow = true
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

	// Sort shifts by start time
	for i := 0; i < len(resp)-1; i++ {
		for j := i + 1; j < len(resp); j++ {
			timeI, errI := time.Parse(time.RFC3339, resp[i].StartTime)
			timeJ, errJ := time.Parse(time.RFC3339, resp[j].StartTime)
			if errI == nil && errJ == nil && timeJ.Before(timeI) {
				resp[i], resp[j] = resp[j], resp[i]
			}
		}
	}

	// if time window provided, filter results
	if filterByWindow {
		filtered := make([]services.AssignedShiftResponse, 0, len(resp))
		for _, r := range resp {
			if r.StartTime == "" || r.EndTime == "" {
				continue
			}
			st, err1 := time.Parse(time.RFC3339, r.StartTime)
			et, err2 := time.Parse(time.RFC3339, r.EndTime)
			if err1 != nil || err2 != nil {
				continue
			}
			// include shift if it overlaps the provided window
			if et.Before(startTime) || st.After(endTime) {
				continue
			}
			filtered = append(filtered, r)
		}
		c.JSON(http.StatusOK, filtered)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateAvailability allows a worker to set their availability for a specific day of the week
func CreateAvailability(c *gin.Context) {
	type Req struct {
		DayOfWeek *int    `json:"day_of_week" binding:"required,min=0,max=6"`
		StartTime string `json:"start_time" binding:"required"`
		EndTime   string `json:"end_time" binding:"required"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateAvailability: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	empParam := c.Param("employee_id")
	if empParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id path parameter is required"})
		return
	}

	empID64, err := strconv.ParseUint(empParam, 10, 64)
	if err != nil {
		log.Printf("CreateAvailability: invalid employee_id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}

	if req.DayOfWeek == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "day_of_week is required"})
		return
	}

	log.Printf("Creating availability for worker %d: day=%d, start=%s, end=%s", empID64, *req.DayOfWeek, req.StartTime, req.EndTime)

	availability, err := services.CreateWorkerAvailability(uint(empID64), *req.DayOfWeek, req.StartTime, req.EndTime)
	if err != nil {
		log.Printf("CreateAvailability: service error: %v", err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "worker not found"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create availability"})
		return
	}

	log.Printf("Availability created successfully, id=%d", availability.ID)
	c.JSON(http.StatusCreated, gin.H{
		"id":          availability.ID,
		"worker_id":   availability.WorkerID,
		"day_of_week": availability.DayOfWeek,
		"start_time":  availability.StartTime,
		"end_time":    availability.EndTime,
		"created_at":  availability.CreatedAt.Format(time.RFC3339),
	})
}

// GetReleasedShiftsForWorkerCompany returns all released shifts in the authenticated worker's company.
func GetReleasedShiftsForWorkerCompany(c *gin.Context) {
	workerID, err := services.GetAuthenticatedEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}

	svc := services.WorkerService{}
	shifts, err := svc.GetReleasedShiftsByEmployeeCompany(workerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "worker not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch released shifts"})
		return
	}

	c.JSON(http.StatusOK, shifts)
}
