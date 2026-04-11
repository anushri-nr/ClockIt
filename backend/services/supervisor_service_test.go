package services

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"clockit/backend/database"
	"clockit/backend/models"
	"clockit/backend/repository"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupServiceTestDB(t *testing.T) {
	t.Helper()

	// create a unique temp directory for this test and place the sqlite file inside it
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	// initialize the database at this per-test path
	if err := database.Init(path); err != nil {
		t.Fatalf("failed to init database: %v", err)
	}

	// ensure the underlying sql.DB is closed and DB pointer cleared after the test
	t.Cleanup(func() {
		if database.DB != nil {
			sqlDB, err := database.DB.DB()
			if err == nil && sqlDB != nil {
				_ = sqlDB.Close()
			}
			database.DB = nil
		}
		// try removing file just in case (TempDir will remove directory tree afterwards)
		_ = os.Remove(path)
	})
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

func TestSupervisorService_GetWorkersWithOvertimeHours_Success(t *testing.T) {
	setupServiceTestDB(t)

	// create company and supervisor
	comp := models.Company{Name: "OTCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	sup, err := RegisterEmployee("SupOT", "supot@example.com", "supass", "addr", "222", comp.ID, 20.0, models.RoleSupervisor)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	// create worker
	worker, err := RegisterEmployee("WOT", "wot@example.com", "pass1234", "addr", "333", comp.ID, 9.0, models.RoleWorker)
	if err != nil {
		t.Fatalf("RegisterEmployee worker failed: %v", err)
	}

	now := time.Now()
	weekday := int(now.Weekday())
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, now.Location())

	// create 3 shifts of 8 hours within the week
	for i := 1; i <= 3; i++ {
		start := weekStart.Add(time.Duration(i) * 24 * time.Hour)
		end := start.Add(8 * time.Hour)
		sft := models.Shift{StartTime: start, EndTime: end, CreatedBy: sup.ID, CreatedAt: time.Now()}
		if err := database.DB.Create(&sft).Error; err != nil {
			t.Fatalf("failed to create shift: %v", err)
		}
		sa := models.ShiftAssignment{ShiftID: sft.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned}
		if err := database.DB.Create(&sa).Error; err != nil {
			t.Fatalf("failed to create assignment: %v", err)
		}
	}

	svc := SupervisorService{}
	resp, err := svc.GetWorkersWithOvertimeHours(sup.ID, weekStart)
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, worker.ID, resp[0].ID)
	require.Greater(t, resp[0].TotalHours, 20.0)
}

func TestSupervisorService_GetWorkersWithOvertimeHours_NotFound(t *testing.T) {
	setupServiceTestDB(t)

	svc := SupervisorService{}
	now := time.Now()
	weekday := int(now.Weekday())
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, now.Location())

	_, err := svc.GetWorkersWithOvertimeHours(99999, weekStart)
	require.Error(t, err)
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestRejectShiftRequest_Success(t *testing.T) {
	setupServiceTestDB(t) // use your existing service DB helper

	company := models.Company{Name: "Reject Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	supervisor := models.Employee{
		Name:      "Sup",
		Email:     "sup.reject@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: company.ID,
		Wage:      40,
	}
	worker := models.Employee{
		Name:      "Worker",
		Email:     "worker.reject@test.local",
		Role:      models.RoleWorker,
		CompanyID: company.ID,
		Wage:      20,
	}
	require.NoError(t, database.DB.Create(&supervisor).Error)
	require.NoError(t, database.DB.Create(&worker).Error)

	now := time.Now().UTC()
	shift := models.Shift{
		StartTime: now.Add(2 * time.Hour),
		EndTime:   now.Add(10 * time.Hour),
		CreatedBy: supervisor.ID,
	}
	require.NoError(t, database.DB.Create(&shift).Error)

	assignment := models.ShiftAssignment{
		ShiftID:    shift.ID,
		EmployeeID: worker.ID,
		AssigneeID: supervisor.ID,
		Status:     models.StatusRequested,
		AssignedAt: now,
	}
	require.NoError(t, database.DB.Create(&assignment).Error)

	svc := SupervisorService{}
	out, err := svc.RejectShiftRequest(supervisor.ID, shift.ID)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, models.StatusReleased, out.Status)
	require.Equal(t, supervisor.ID, out.AssigneeID)

	var dbRow models.ShiftAssignment
	require.NoError(t, database.DB.Where("shift_id = ?", shift.ID).First(&dbRow).Error)
	require.Equal(t, models.StatusReleased, dbRow.Status)
	require.Equal(t, supervisor.ID, dbRow.AssigneeID)
}

func TestRejectShiftRequest_NotFound_WrongCompany(t *testing.T) {
	setupServiceTestDB(t) // use your existing service DB helper

	companyA := models.Company{Name: "A"}
	companyB := models.Company{Name: "B"}
	require.NoError(t, database.DB.Create(&companyA).Error)
	require.NoError(t, database.DB.Create(&companyB).Error)

	supA := models.Employee{
		Name:      "SupA",
		Email:     "supa.reject@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: companyA.ID,
		Wage:      40,
	}
	supB := models.Employee{
		Name:      "SupB",
		Email:     "supb.reject@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: companyB.ID,
		Wage:      40,
	}
	workerB := models.Employee{
		Name:      "WorkerB",
		Email:     "workerb.reject@test.local",
		Role:      models.RoleWorker,
		CompanyID: companyB.ID,
		Wage:      20,
	}
	require.NoError(t, database.DB.Create(&supA).Error)
	require.NoError(t, database.DB.Create(&supB).Error)
	require.NoError(t, database.DB.Create(&workerB).Error)

	now := time.Now().UTC()
	shiftB := models.Shift{
		StartTime: now.Add(3 * time.Hour),
		EndTime:   now.Add(11 * time.Hour),
		CreatedBy: supB.ID,
	}
	require.NoError(t, database.DB.Create(&shiftB).Error)

	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID:    shiftB.ID,
		EmployeeID: workerB.ID,
		AssigneeID: supB.ID,
		Status:     models.StatusRequested,
		AssignedAt: now,
	}).Error)

	svc := SupervisorService{}
	out, err := svc.RejectShiftRequest(supA.ID, shiftB.ID)
	require.Error(t, err)
	require.Nil(t, out)
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
