package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// WorkerAvailabilityResponse format
type WorkerAvailabilityResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	PhoneNo   string `json:"phone_no"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// RegisterSupervisor creates a new supervisor
func RegisterSupervisor(name, email, password, address, phone string, companyID uint, wage float64) (*models.Employee, error) {
	// Check if email already exists
	log.Printf("RegisterSupervisor: checking if email exists: %s", email)
	var existing models.Employee
	if database.DB.Where("email = ?", email).First(&existing).RowsAffected == 1 {
		log.Printf("RegisterSupervisor: email already exists: %s", email)
		return nil, errors.New("employee with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("RegisterSupervisor: password hashing failed: %v", err)
		return nil, err
	}

	employee := models.Employee{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      models.RoleSupervisor,
		Address:   address,
		PhoneNo:   phone,
		CompanyID: companyID,
		Wage:      wage,
	}

	if err := database.DB.Create(&employee).Error; err != nil {
		log.Printf("RegisterSupervisor: DB insert failed: %v", err)
		return nil, err
	}

	log.Printf("RegisterSupervisor: supervisor created, id=%d", employee.ID)
	return &employee, nil
}