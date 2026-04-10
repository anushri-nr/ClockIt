package services

import (
	"clockit/backend/database"
	"clockit/backend/models"
	"clockit/backend/repository"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
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

type CreatedShiftsResponse struct {
	ID        uint   `json:"id"`
	ShiftID   uint   `json:"shift_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	CreatedAt string `json:"created_at"`

	AssignedBy uint   `json:"assigned_by,omitempty"`
	AssignedTo string `json:"assigned_to,omitempty"`
	Status     string `json:"status,omitempty"`
}

type SupervisorService struct{}

// GetShiftsByCompany returns all shifts for a supervisor's company, optionally filtered by status
func (s *SupervisorService) GetShiftsByCompany(supervisorID uint, statusFilter string) ([]CreatedShiftsResponse, error) {
	var sup models.Employee
	if err := database.DB.Where("id = ? AND role = ?", supervisorID, models.RoleSupervisor).First(&sup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("GetShiftsByCompany: supervisor not found: %v", err)
			return []CreatedShiftsResponse{}, gorm.ErrRecordNotFound
		}

		log.Printf("GetShiftsByCompany: DB error validating supervisor: %v", err)

		return []CreatedShiftsResponse{}, err
	}

	// Base Query: Filters by Company ID
	query := database.DB.Model(&models.Shift{}).
		Joins("JOIN employees ON employees.id = shifts.created_by").
		Where("employees.company_id = ?", sup.CompanyID).
		Preload("Creator")

	// Only Group By when we actually execute a JOIN on assignments
	if statusFilter != "" {
		query = query.Joins("LEFT JOIN shift_assignments ON shift_assignments.shift_id = shifts.id").
			Group("shifts.id")

		// If Unassigned, only check for IS NULL since it's not a saved DB value
		if statusFilter == string(models.StatusUnassigned) {
			query = query.Where("shift_assignments.id IS NULL")
			// don't preload assignments when unassigned
		} else {
			// Filter shifts to those that have an assignment with the requested status
			query = query.Where("shift_assignments.status = ?", statusFilter)
			// Preload only assignments that match the filter, order most recent first
			query = query.Preload("Assignments", func(db *gorm.DB) *gorm.DB {
				return db.Where("status = ?", statusFilter).Order("assigned_at DESC")
			}).Preload("Assignments.Employee").Preload("Assignments.Assignee")
		}
	} else {
		// No status filter: preload assignments ordered by assigned_at desc so first is latest
		query = query.Preload("Assignments", func(db *gorm.DB) *gorm.DB {
			return db.Order("assigned_at DESC")
		}).Preload("Assignments.Employee").Preload("Assignments.Assignee")
	}

	var shifts []models.Shift
	if err := query.Find(&shifts).Error; err != nil {
		log.Printf("GetShiftsByCompany: DB query failed: %v", err)
		return []CreatedShiftsResponse{}, err
	}

	resp := make([]CreatedShiftsResponse, 0, len(shifts))
	for _, sft := range shifts {
		var assignedBy uint
		var assignedTo string

		// Use the new Enum as the default string
		status := string(models.StatusUnassigned)

		if len(sft.Assignments) > 0 {
			// After conditional preload the first assignment is the most relevant (matching filter or latest)
			a := sft.Assignments[0]
			assignedBy = a.AssigneeID
			status = string(a.Status)
			if a.Employee.ID != 0 {
				assignedTo = a.Employee.Name
			}
		}

		sc := CreatedShiftsResponse{
			ID:         sft.ID,
			ShiftID:    sft.ID,
			StartTime:  sft.StartTime.Format(time.RFC3339),
			EndTime:    sft.EndTime.Format(time.RFC3339),
			CreatedAt:  sft.CreatedAt.Format(time.RFC3339),
			AssignedBy: assignedBy,
			AssignedTo: assignedTo,
			Status:     status,
		}
		resp = append(resp, sc)
	}

	return resp, nil
}

type WorkerOvertimeResponse struct {
	ID         uint    `json:"id"`
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	PhoneNo    string  `json:"phone_no"`
	CompanyID  uint    `json:"company_id"`
	TotalHours float64 `json:"total_hours"`
}

// GetWorkersWithOvertimeHours returns workers whose total assigned shift hours in the week exceed 20 hours.
func (s *SupervisorService) GetWorkersWithOvertimeHours(supervisorID uint, weekStart time.Time) ([]WorkerOvertimeResponse, error) {
	var sup models.Employee
	if err := database.DB.Where("id = ? AND role = ?", supervisorID, models.RoleSupervisor).First(&sup).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []WorkerOvertimeResponse{}, gorm.ErrRecordNotFound
		}
		return []WorkerOvertimeResponse{}, err
	}

	weekEnd := weekStart.AddDate(0, 0, 7)

	sql := `SELECT e.id, e.name, e.email, e.phone_no, e.company_id,
		SUM((julianday(shifts.end_time) - julianday(shifts.start_time)) * 24.0) as total_hours
		FROM employees e
		JOIN shift_assignments sa ON sa.employee_id = e.id
		JOIN shifts ON sa.shift_id = shifts.id
		WHERE e.company_id = ? AND shifts.start_time >= ? AND shifts.start_time < ?
		GROUP BY e.id
		HAVING total_hours > ?`

	var out []WorkerOvertimeResponse
	if err := database.DB.Raw(sql, sup.CompanyID, weekStart.Format(time.RFC3339), weekEnd.Format(time.RFC3339), 20).Scan(&out).Error; err != nil {
		return []WorkerOvertimeResponse{}, err
	}

	return out, nil
}
