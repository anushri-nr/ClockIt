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

func TestSupervisorService_GetShiftsByCompany_UnassignedFilter(t *testing.T) {
	setupServiceTestDB(t)

	comp := models.Company{Name: "UnassignedSvcCo"}
	require.NoError(t, database.DB.Create(&comp).Error)

	sup := models.Employee{Name: "Sup", Email: "sup-unassigned@example.com", Password: "x", Role: models.RoleSupervisor, CompanyID: comp.ID}
	require.NoError(t, database.DB.Create(&sup).Error)

	start := time.Now().Add(48 * time.Hour)
	assignedShift := models.Shift{StartTime: start, EndTime: start.Add(8 * time.Hour), CreatedBy: sup.ID, CreatedAt: time.Now()}
	unassignedShift := models.Shift{StartTime: start.Add(24 * time.Hour), EndTime: start.Add(32 * time.Hour), CreatedBy: sup.ID, CreatedAt: time.Now()}
	require.NoError(t, database.DB.Create(&assignedShift).Error)
	require.NoError(t, database.DB.Create(&unassignedShift).Error)

	worker := models.Employee{Name: "Worker", Email: "worker-unassigned@example.com", Password: "x", Role: models.RoleWorker, CompanyID: comp.ID}
	require.NoError(t, database.DB.Create(&worker).Error)
	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID: assignedShift.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned,
	}).Error)

	svc := SupervisorService{}
	resp, err := svc.GetShiftsByCompany(sup.ID, string(models.StatusUnassigned))
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, unassignedShift.ID, resp[0].ShiftID)
	require.Equal(t, string(models.StatusUnassigned), resp[0].Status)
}

func TestSupervisorService_GetShiftsByCompany_NoFilterUsesLatestAssignment(t *testing.T) {
	setupServiceTestDB(t)

	comp := models.Company{Name: "LatestSvcCo"}
	require.NoError(t, database.DB.Create(&comp).Error)

	sup := models.Employee{Name: "Sup", Email: "sup-latest@example.com", Password: "x", Role: models.RoleSupervisor, CompanyID: comp.ID}
	oldWorker := models.Employee{Name: "Old Worker", Email: "old-worker-latest@example.com", Password: "x", Role: models.RoleWorker, CompanyID: comp.ID}
	worker := models.Employee{Name: "Worker", Email: "worker-latest@example.com", Password: "x", Role: models.RoleWorker, CompanyID: comp.ID}
	require.NoError(t, database.DB.Create(&sup).Error)
	require.NoError(t, database.DB.Create(&oldWorker).Error)
	require.NoError(t, database.DB.Create(&worker).Error)

	start := time.Now().Add(72 * time.Hour)
	shift := models.Shift{StartTime: start, EndTime: start.Add(8 * time.Hour), CreatedBy: sup.ID, CreatedAt: time.Now()}
	require.NoError(t, database.DB.Create(&shift).Error)
	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID: shift.ID, EmployeeID: oldWorker.ID, AssigneeID: sup.ID, AssignedAt: time.Now().Add(-1 * time.Hour), Status: models.StatusReleased,
	}).Error)
	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID: shift.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned,
	}).Error)

	svc := SupervisorService{}
	resp, err := svc.GetShiftsByCompany(sup.ID, "")
	require.NoError(t, err)
	require.Len(t, resp, 1)
	require.Equal(t, string(models.StatusAssigned), resp[0].Status)
	require.Equal(t, worker.Name, resp[0].AssignedTo)
	require.Equal(t, worker.ID, resp[0].AssignedToID)
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

func TestSupervisorService_GetWorkersWithOvertimeHours_ExcludesWorkersAtOrBelowThreshold(t *testing.T) {
	setupServiceTestDB(t)

	comp := models.Company{Name: "NoOTCo"}
	require.NoError(t, database.DB.Create(&comp).Error)

	sup, err := RegisterEmployee("SupNoOT", "supnoot@example.com", "supass", "addr", "222", comp.ID, 20.0, models.RoleSupervisor)
	require.NoError(t, err)

	worker, err := RegisterEmployee("WNoOT", "wnoot@example.com", "pass1234", "addr", "333", comp.ID, 9.0, models.RoleWorker)
	require.NoError(t, err)

	now := time.Now()
	weekday := int(now.Weekday())
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, now.Location())

	start := weekStart.Add(24 * time.Hour)
	sft := models.Shift{StartTime: start, EndTime: start.Add(8 * time.Hour), CreatedBy: sup.ID, CreatedAt: time.Now()}
	require.NoError(t, database.DB.Create(&sft).Error)
	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID: sft.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned,
	}).Error)

	svc := SupervisorService{}
	resp, err := svc.GetWorkersWithOvertimeHours(sup.ID, weekStart)
	require.NoError(t, err)
	require.Len(t, resp, 0)
}
