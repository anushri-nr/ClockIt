package services

import (
	"clockit/backend/repository"
	"log"
	"time"
)

// WorkerAvailabilityResponse format
type WorkerAvailabilityResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNo     string `json:"phone_no"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	CompanyID   uint   `json:"company_id"`
	CompanyName string `json:"company_name"`
}

type WorkerService struct {
	Repo *repository.WorkerRepository
}

func (s *WorkerService) GetWorkersAvailable(date time.Time, companyID uint) ([]WorkerAvailabilityResponse, error) {
	availabilities, err := s.Repo.FindAvailableWorkers(date, companyID)
	if err != nil {
		log.Printf("GetWorkersAvailableForDate: fetching workers failed: %v", err)
		return []WorkerAvailabilityResponse{}, err
	}

	response := []WorkerAvailabilityResponse{}
	for _, a := range availabilities {
		response = append(response, WorkerAvailabilityResponse{
			ID:          a.Worker.ID,
			Name:        a.Worker.Name,
			Email:       a.Worker.Email,
			PhoneNo:     a.Worker.PhoneNo,
			StartTime:   a.StartTime,
			EndTime:     a.EndTime,
			CompanyID:   a.Worker.CompanyID,
			CompanyName: a.Worker.Company.Name,
		})
	}

	log.Printf("GetWorkersAvailable: returning %d worker availability entries", len(response))
	return response, nil
}
