package controllers

import (
	"clockit/backend/services"
	"log"
	"net/http"

	"clockit/backend/models"

	"github.com/gin-gonic/gin"
)

func WorkerRegister(c *gin.Context) {
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
		log.Printf("WorkerRegister: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	if req.Wage < 0 {
		log.Printf("WorkerRegister: negative wage: %f", req.Wage)
		c.JSON(http.StatusBadRequest, gin.H{"error": "wage must be non-negative"})
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
