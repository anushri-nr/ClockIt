package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// RegisterEmployee creates a new employee (worker or supervisor) in the database
func RegisterEmployee(name, email, password, address, phone string, companyID uint, wage float64, role models.EmployeeRole) (*models.Employee, error) {
	// Check if email already exists
	log.Printf("RegisterEmployee: checking if email exists: %s", email)
	var existing models.Employee
	if database.DB.Where("email = ?", email).First(&existing).RowsAffected == 1 {
		log.Printf("RegisterEmployee: email already exists: %s", email)
		return nil, errors.New("employee with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("RegisterEmployee: password hashing failed: %v", err)
		return nil, err
	}

	employee := models.Employee{
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		Role:      role,
		Address:   address,
		PhoneNo:   phone,
		CompanyID: companyID,
		Wage:      wage,
	}

	if err := database.DB.Create(&employee).Error; err != nil {
		log.Printf("RegisterEmployee: DB insert failed: %v", err)
		return nil, err
	}

	log.Printf("RegisterEmployee: employee created, id=%d", employee.ID)
	return &employee, nil
}
