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
	err := database.DB.Where("email = ?", email).First(&existing).Error
	if err == nil {
		// Record found, email already exists
		log.Printf("RegisterEmployee: email already exists: %s", email)
		return nil, errors.New("employee with this email already exists")
	}

	// If error is not "record not found", it's a database error
	if err.Error() != "record not found" {
		log.Printf("RegisterEmployee: database error while checking email: %v", err)
		return nil, err
	}

	// Validate wage
	if wage < 0 {
		log.Printf("RegisterEmployee: invalid wage: %f", wage)
		return nil, errors.New("wage must be non-negative")
	}

	//Validate company exists
	var company models.Company
	if err := database.DB.First(&company, companyID).Error; err != nil {
		log.Printf("RegisterEmployee: company not found: %v", err)
		return nil, errors.New("company not found")
	}

	// Validate role
	if role != models.RoleSupervisor && role != models.RoleWorker {
		log.Printf("RegisterEmployee: invalid role: %s", role)
		return nil, errors.New("invalid role")
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
