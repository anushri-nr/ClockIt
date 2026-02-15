package repository

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"log"
	"time"
)

type WorkerRepository struct{}

// FindAvailableWorkers returns all worker availabilities for a given date and company
func (r *WorkerRepository) FindAvailableWorkers(date time.Time, companyID uint) ([]models.WorkerAvailability, error) {
	dayOfWeek := int(date.Weekday())
	var availabilities []models.WorkerAvailability

	err := database.DB.
		Model(&models.WorkerAvailability{}).
		Joins("JOIN employees ON employees.id = worker_availabilities.worker_id").
		Joins("JOIN companies ON companies.id = employees.company_id").
		Where("worker_availabilities.day_of_week = ?", dayOfWeek).
		Where("employees.company_id = ?", companyID).
		Where("employees.role = ?", models.RoleWorker).
		Preload("Worker").
		Preload("Worker.Company").
		Find(&availabilities).Error

	if err != nil {
		log.Printf("GetWorkersAvailableForDate: DB query error: %v", err)
		return []models.WorkerAvailability{}, err
	}

	return availabilities, nil
}
