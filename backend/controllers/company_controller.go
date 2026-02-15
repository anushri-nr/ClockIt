package controllers

import (
	"clockit/backend/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateCompany handles company creation
func CreateCompany(c *gin.Context) {
	type Req struct {
		Name    string `json:"name" binding:"required"`
		Address string `json:"address"`
	}

	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateCompany: binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": GetValidationErrors(err)})
		return
	}

	company, err := services.CreateCompany(req.Name, req.Address)
	if err != nil {
		log.Printf("CreateCompany: failed to create company: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         company.ID,
		"name":       company.Name,
		"address":    company.Address,
		"created_at": company.CreatedAt,
		"updated_at": company.UpdatedAt,
	})
}
