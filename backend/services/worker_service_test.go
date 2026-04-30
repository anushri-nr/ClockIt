package services

import (
	"errors"
	"testing"
	"time"

	"clockit/backend/database"
	"clockit/backend/models"
	"clockit/backend/repository"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWorkerServiceDB(t *testing.T) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	database.DB = db
	require.NoError(t, database.DB.AutoMigrate(
		&models.Company{},
		&models.Employee{},
		&models.WorkerAvailability{},
		&models.Shift{},
		&models.ShiftAssignment{},
	))

	t.Cleanup(func() {
		if database.DB != nil {
			if sqlDB, err := database.DB.DB(); err == nil {
				_ = sqlDB.Close()
			}
			database.DB = nil
		}
	})
}

func TestCreateWorkerAvailability(t *testing.T) {
	t.Run("creates availability for valid worker", func(t *testing.T) {
		setupWorkerServiceDB(t)
		worker := models.Employee{
			ID:    1,
			Email: "worker1@test.local",
			Role:  models.RoleWorker,
			Wage:  20,
		}
		require.NoError(t, database.DB.Create(&worker).Error)

		got, err := CreateWorkerAvailability(1, 1, "09:00:00", "17:00:00")
		require.NoError(t, err)
		require.NotNil(t, got)
		require.Equal(t, uint(1), got.WorkerID)
		require.Equal(t, 1, got.DayOfWeek)
		require.Equal(t, "09:00:00", got.StartTime)
		require.Equal(t, "17:00:00", got.EndTime)
	})

	t.Run("returns not found for non-worker", func(t *testing.T) {
		setupWorkerServiceDB(t)
		supervisor := models.Employee{
			ID:    2,
			Email: "supervisor2@test.local",
			Role:  models.RoleSupervisor,
			Wage:  30,
		}
		require.NoError(t, database.DB.Create(&supervisor).Error)

		got, err := CreateWorkerAvailability(2, 2, "10:00:00", "18:00:00")
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("returns not found for missing worker", func(t *testing.T) {
		setupWorkerServiceDB(t)
		got, err := CreateWorkerAvailability(999, 3, "08:00:00", "16:00:00")
		require.Error(t, err)
		require.Nil(t, got)
		require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})
}

func TestGetAssignedShifts_Success(t *testing.T) {
	setupWorkerServiceDB(t)

	company := models.Company{Name: "Assigned Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	supervisor := models.Employee{
		Name: "Supervisor", Email: "supervisor.assigned@test.local", Role: models.RoleSupervisor, CompanyID: company.ID, Wage: 40,
	}
	worker := models.Employee{
		Name: "Worker", Email: "worker2@test.local", Role: models.RoleWorker, CompanyID: company.ID, Wage: 20,
	}
	require.NoError(t, database.DB.Create(&supervisor).Error)
	require.NoError(t, database.DB.Create(&worker).Error)

	start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	shift := models.Shift{StartTime: start, EndTime: start.Add(6 * time.Hour), CreatedBy: supervisor.ID}
	require.NoError(t, database.DB.Create(&shift).Error)

	assignedAt := time.Now().UTC().Truncate(time.Second)
	assignment := models.ShiftAssignment{
		ShiftID: shift.ID, EmployeeID: worker.ID, AssigneeID: supervisor.ID, Status: models.StatusAssigned, AssignedAt: assignedAt,
	}
	require.NoError(t, database.DB.Create(&assignment).Error)

	svc := &WorkerService{Repo: &repository.WorkerRepository{}}
	got, err := svc.GetAssignedShifts(worker.ID)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, assignment.ID, got[0].ID)
	require.Equal(t, shift.ID, got[0].ShiftID)
	require.Equal(t, supervisor.ID, got[0].AssignedBy)
	require.Equal(t, string(models.StatusAssigned), got[0].Status)
	require.Equal(t, 6.0, got[0].DurationHours)
	require.Equal(t, 120.0, got[0].Earnings)
	require.Equal(t, start.Format(time.RFC3339), got[0].StartTime)
	require.Equal(t, start.Add(6*time.Hour).Format(time.RFC3339), got[0].EndTime)
	require.Equal(t, assignedAt.Format(time.RFC3339), got[0].AssignedAt)
}

func TestGetAssignedShifts(t *testing.T) {
	setupWorkerServiceDB(t)

	svc := &WorkerService{} // repo is not called in these validation-failure paths

	t.Run("returns not found when employee does not exist", func(t *testing.T) {
		got, err := svc.GetAssignedShifts(999)
		require.Error(t, err)
		require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		require.Len(t, got, 0)
	})

	t.Run("returns not found when employee is not worker", func(t *testing.T) {
		supervisor := models.Employee{
			ID:   10,
			Role: models.RoleSupervisor,
			Wage: 40,
		}
		require.NoError(t, database.DB.Create(&supervisor).Error)

		got, err := svc.GetAssignedShifts(10)
		require.Error(t, err)
		require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		require.Len(t, got, 0)
	})
}

func TestGetReleasedShiftsByEmployeeCompany_Success(t *testing.T) {
	setupWorkerServiceDB(t)

	companyA := models.Company{Name: "A"}
	companyB := models.Company{Name: "B"}
	require.NoError(t, database.DB.Create(&companyA).Error)
	require.NoError(t, database.DB.Create(&companyB).Error)

	worker := models.Employee{
		Name: "Worker", Email: "worker.rel@test.local", Role: models.RoleWorker, CompanyID: companyA.ID, Wage: 20,
	}
	supA := models.Employee{
		Name: "SupA", Email: "supa.rel@test.local", Role: models.RoleSupervisor, CompanyID: companyA.ID, Wage: 40,
	}
	supB := models.Employee{
		Name: "SupB", Email: "supb.rel@test.local", Role: models.RoleSupervisor, CompanyID: companyB.ID, Wage: 40,
	}
	require.NoError(t, database.DB.Create(&worker).Error)
	require.NoError(t, database.DB.Create(&supA).Error)
	require.NoError(t, database.DB.Create(&supB).Error)

	now := time.Now().UTC()
	shiftKeep := models.Shift{StartTime: now.Add(1 * time.Hour), EndTime: now.Add(9 * time.Hour), CreatedBy: supA.ID}
	shiftDropStatus := models.Shift{StartTime: now.Add(2 * time.Hour), EndTime: now.Add(10 * time.Hour), CreatedBy: supA.ID}
	shiftDropCompany := models.Shift{StartTime: now.Add(3 * time.Hour), EndTime: now.Add(11 * time.Hour), CreatedBy: supB.ID}
	require.NoError(t, database.DB.Create(&shiftKeep).Error)
	require.NoError(t, database.DB.Create(&shiftDropStatus).Error)
	require.NoError(t, database.DB.Create(&shiftDropCompany).Error)

	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID: shiftKeep.ID, EmployeeID: worker.ID, Status: models.StatusReleased, AssignedAt: now,
	}).Error)
	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID: shiftDropStatus.ID, EmployeeID: worker.ID, Status: models.StatusAssigned, AssignedAt: now,
	}).Error)
	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID: shiftDropCompany.ID, EmployeeID: worker.ID, Status: models.StatusReleased, AssignedAt: now,
	}).Error)

	svc := WorkerService{}
	got, err := svc.GetReleasedShiftsByEmployeeCompany(worker.ID)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, shiftKeep.ID, got[0].ID)
}

func TestGetReleasedShiftsByEmployeeCompany_EmployeeNotWorker(t *testing.T) {
	setupWorkerServiceDB(t)

	company := models.Company{Name: "Released Wrong Role Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	supervisor := models.Employee{
		Name: "SupWrongRole", Email: "sup.wrongrole.rel@test.local", Role: models.RoleSupervisor, CompanyID: company.ID, Wage: 40,
	}
	require.NoError(t, database.DB.Create(&supervisor).Error)

	svc := WorkerService{}
	got, err := svc.GetReleasedShiftsByEmployeeCompany(supervisor.ID)
	require.Error(t, err)
	require.Len(t, got, 0)
	require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestGetReleasedShiftsByEmployeeCompany_Empty(t *testing.T) {
	setupWorkerServiceDB(t)

	company := models.Company{Name: "Released Empty Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	worker := models.Employee{
		Name: "WorkerEmpty", Email: "worker.empty.rel@test.local", Role: models.RoleWorker, CompanyID: company.ID, Wage: 20,
	}
	require.NoError(t, database.DB.Create(&worker).Error)

	svc := WorkerService{}
	got, err := svc.GetReleasedShiftsByEmployeeCompany(worker.ID)
	require.NoError(t, err)
	require.Len(t, got, 0)
}

func TestRequestReleasedShift_Success(t *testing.T) {
	setupWorkerServiceDB(t)

	company := models.Company{Name: "Req Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	worker := models.Employee{
		Name: "WorkerReq", Email: "worker.req@test.local", Role: models.RoleWorker, CompanyID: company.ID, Wage: 22,
	}
	supervisor := models.Employee{
		Name: "SupReq", Email: "sup.req@test.local", Role: models.RoleSupervisor, CompanyID: company.ID, Wage: 45,
	}
	require.NoError(t, database.DB.Create(&worker).Error)
	require.NoError(t, database.DB.Create(&supervisor).Error)

	now := time.Now().UTC()
	shift := models.Shift{
		StartTime: now.Add(24 * time.Hour),
		EndTime:   now.Add(32 * time.Hour),
		CreatedBy: supervisor.ID,
	}
	require.NoError(t, database.DB.Create(&shift).Error)

	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID:    shift.ID,
		EmployeeID: worker.ID,
		AssigneeID: supervisor.ID,
		Status:     models.StatusReleased,
		AssignedAt: now,
	}).Error)

	out, err := RequestReleasedShift(worker.ID, shift.ID)
	require.NoError(t, err)
	require.NotNil(t, out)
	require.Equal(t, models.StatusRequested, out.Status)
	require.Equal(t, worker.ID, out.EmployeeID)

	var dbRow models.ShiftAssignment
	require.NoError(t, database.DB.Where("shift_id = ?", shift.ID).First(&dbRow).Error)
	require.Equal(t, models.StatusRequested, dbRow.Status)
	require.Equal(t, worker.ID, dbRow.EmployeeID)
}

func TestRequestReleasedShift_NotReleased(t *testing.T) {
	setupWorkerServiceDB(t)

	company := models.Company{Name: "Req Co 2"}
	require.NoError(t, database.DB.Create(&company).Error)

	worker := models.Employee{
		Name: "Worker2", Email: "worker2.req@test.local", Role: models.RoleWorker, CompanyID: company.ID, Wage: 22,
	}
	supervisor := models.Employee{
		Name: "Sup2", Email: "sup2.req@test.local", Role: models.RoleSupervisor, CompanyID: company.ID, Wage: 45,
	}
	require.NoError(t, database.DB.Create(&worker).Error)
	require.NoError(t, database.DB.Create(&supervisor).Error)

	now := time.Now().UTC()
	shift := models.Shift{
		StartTime: now.Add(24 * time.Hour),
		EndTime:   now.Add(32 * time.Hour),
		CreatedBy: supervisor.ID,
	}
	require.NoError(t, database.DB.Create(&shift).Error)

	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID:    shift.ID,
		EmployeeID: worker.ID,
		AssigneeID: supervisor.ID,
		Status:     models.StatusAssigned, // not released
		AssignedAt: now,
	}).Error)

	out, err := RequestReleasedShift(worker.ID, shift.ID)
	require.Error(t, err)
	require.Nil(t, out)
}
