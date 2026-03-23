package services

import (
	"os"
	"testing"
	"time"

	"clockit/backend/database"
	"clockit/backend/models"
	"clockit/backend/repository"

	"github.com/stretchr/testify/require"
)

func setupServiceTestDB(t *testing.T) {
	t.Helper()
	path := "database/test.db"
	_ = os.Remove(path)
	if err := database.Init(path); err != nil {
		t.Fatalf("failed to init database: %v", err)
	}
}

func TestSupervisorService_GetShiftsByCompany_NotFound(t *testing.T) {
	setupServiceTestDB(t)

	svc := SupervisorService{}
	_, err := svc.GetShiftsByCompany(9999, "")
	require.Error(t, err)
}

func TestSupervisorService_GetShiftsByCompany_Success(t *testing.T) {
	setupServiceTestDB(t)

	// create company and supervisor
	comp := models.Company{Name: "SvcCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}
	sup := models.Employee{Name: "SS", Email: "ss@example.com", Password: "x", Role: models.RoleSupervisor, CompanyID: comp.ID}
	if err := database.DB.Create(&sup).Error; err != nil {
		t.Fatalf("failed to create supervisor: %v", err)
	}

	w := models.Employee{Name: "WW", Email: "ww@example.com", Password: "x", Role: models.RoleWorker, CompanyID: comp.ID}
	if err := database.DB.Create(&w).Error; err != nil {
		t.Fatalf("failed to create worker: %v", err)
	}

	start := time.Now().Add(48 * time.Hour)
	end := start.Add(8 * time.Hour)
	sft := models.Shift{StartTime: start, EndTime: end, CreatedBy: sup.ID, CreatedAt: time.Now()}
	if err := database.DB.Create(&sft).Error; err != nil {
		t.Fatalf("failed to create shift: %v", err)
	}

	sa := models.ShiftAssignment{ShiftID: sft.ID, EmployeeID: w.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned}
	if err := database.DB.Create(&sa).Error; err != nil {
		t.Fatalf("failed to create assignment: %v", err)
	}

	svc := SupervisorService{}
	resp, err := svc.GetShiftsByCompany(sup.ID, string(models.StatusAssigned))
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, string(models.StatusAssigned), resp[0].Status)
}

func TestWorkerService_GetWorkersAvailable_Success(t *testing.T) {
	setupServiceTestDB(t)

	// create company and worker
	comp := models.Company{Name: "AvailSvcCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	worker, err := RegisterEmployee("WorkerA", "wa@example.com", "pass1234", "addr", "000", comp.ID, 7.5, models.RoleWorker)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	// create availability for a specific date
	dateStr := "2022-03-23"
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t.Fatalf("failed to parse date: %v", err)
	}

	wa := models.WorkerAvailability{
		WorkerID:  worker.ID,
		DayOfWeek: int(parsed.Weekday()),
		StartTime: "08:00",
		EndTime:   "12:00",
	}
	if err := database.DB.Create(&wa).Error; err != nil {
		t.Fatalf("failed to create worker availability: %v", err)
	}

	ws := WorkerService{Repo: &repository.WorkerRepository{}}
	resp, err := ws.GetWorkersAvailable(parsed, comp.ID)
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, worker.ID, resp[0].ID)
	require.Equal(t, wa.StartTime, resp[0].StartTime)
}

func TestWorkerService_GetWorkersAvailable_NoResults(t *testing.T) {
	setupServiceTestDB(t)

	comp := models.Company{Name: "NoResCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	parsed := time.Now()
	ws := WorkerService{Repo: &repository.WorkerRepository{}}
	resp, err := ws.GetWorkersAvailable(parsed, comp.ID)
	require.NoError(t, err)
	require.Len(t, resp, 0)
}
