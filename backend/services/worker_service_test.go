package services

import (
    "errors"
    "testing"

    "clockit/backend/database"
    "clockit/backend/models"

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
        &models.Employee{},
        &models.WorkerAvailability{},
    ))
}

func TestCreateWorkerAvailability(t *testing.T) {
    t.Run("creates availability for valid worker", func(t *testing.T) {
        setupWorkerServiceDB(t)
        worker := models.Employee{
            ID:   1,
			Email: "worker1@test.local",
            Role: models.RoleWorker,
            Wage: 20,
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
            ID:   2,
			Email: "supervisor2@test.local",
            Role: models.RoleSupervisor,
            Wage: 30,
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

    worker := models.Employee{
        ID:   3,
        Email: "worker2@test.local",
        Role: models.RoleWorker,
        Wage: 20,
    }
    require.NoError(t, database.DB.Create(&worker).Error)
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