package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"log"
	"time"
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

// Returns all workers available on a specific date
func GetWorkersAvailableForDate(date time.Time) ([]WorkerAvailabilityResponse, error) {
	dayOfWeek := int(date.Weekday())
	log.Printf("GetWorkersAvailableForDate: checking availability for weekday=%d", dayOfWeek)

	var availabilities []models.WorkerAvailability
	err := database.DB.Where("day_of_week = ?", dayOfWeek).Find(&availabilities).Error
	if err != nil {
		log.Printf("GetWorkersAvailableForDate: DB query error: %v", err)
		return nil, err
	}

	if len(availabilities) == 0 {
		return []WorkerAvailabilityResponse{}, nil
	}

	workerIDs := make([]uint, len(availabilities))
	for i, a := range availabilities {
		workerIDs[i] = a.WorkerID
	}

	var workers []models.Employee
	err = database.DB.Where("id IN ? AND role = ?", workerIDs, models.RoleWorker).Find(&workers).Error
	if err != nil {
		log.Printf("GetWorkersAvailableForDate: fetching workers failed: %v", err)
		return nil, err
	}

	workerMap := make(map[uint]models.Employee)
	for _, w := range workers {
		workerMap[w.ID] = w
	}

	var response []WorkerAvailabilityResponse
	for _, a := range availabilities {
		if worker, ok := workerMap[a.WorkerID]; ok {
			response = append(response, WorkerAvailabilityResponse{
				ID:        worker.ID,
				Name:      worker.Name,
				Email:     worker.Email,
				PhoneNo:   worker.PhoneNo,
				StartTime: a.StartTime,
				EndTime:   a.EndTime,
			})
		}
	}
	log.Printf("GetWorkersAvailableForDate: returning %d worker availability entries", len(response))
	return response, nil
}
