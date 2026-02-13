package controllers

import (
	"clockit/backend/services"
	"log"
	"net/http"

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

	employee, err := services.RegisterSupervisor(
		req.Name,
		req.Email,
		req.Password,
		req.Address,
		req.PhoneNo,
		req.CompanyID,
		req.Wage,
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
