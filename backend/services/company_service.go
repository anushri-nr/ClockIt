package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"errors"
	"log"
)

// CreateCompany creates a new company record
func CreateCompany(name, address string) (*models.Company, error) {
	if name == "" {
		return nil, errors.New("company name is required")
	}

	company := models.Company{
		Name:    name,
		Address: address,
	}

	if err := database.DB.Create(&company).Error; err != nil {
		log.Printf("CreateCompany: DB insert failed: %v", err)
		return nil, err
	}

	log.Printf("CreateCompany: company created, id=%d", company.ID)
	return &company, nil
}
