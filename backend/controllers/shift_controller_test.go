package controllers

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"clockit/backend/database"
	"clockit/backend/models"

	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupShiftControllerTestDB(t *testing.T) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "shift_controller_test.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err)

	database.DB = db

	err = database.DB.AutoMigrate(
		&models.Company{},
		&models.Employee{},
		&models.Shift{},
		&models.ShiftAssignment{},
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		if database.DB != nil {
			if sqlDB, err := database.DB.DB(); err == nil {
				_ = sqlDB.Close()
			}
			database.DB = nil
		}
	})
}

func seedCompanyAndUsersForShiftTests(t *testing.T) {
	t.Helper()

	company := models.Company{Name: "Shift Test Company"}
	require.NoError(t, database.DB.Create(&company).Error)

	creator := models.Employee{
		ID:        1,
		Name:      "Supervisor Creator",
		Email:     "creator@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: company.ID,
		Wage:      30,
	}
	require.NoError(t, database.DB.Create(&creator).Error)

	worker := models.Employee{
		ID:        2,
		Name:      "Worker One",
		Email:     "worker@test.local",
		Role:      models.RoleWorker,
		CompanyID: company.ID,
		Wage:      20,
	}
	require.NoError(t, database.DB.Create(&worker).Error)

	assigner := models.Employee{
		ID:        3,
		Name:      "Supervisor Assigner",
		Email:     "assigner@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: company.ID,
		Wage:      35,
	}
	require.NoError(t, database.DB.Create(&assigner).Error)
}

func seedShiftForAssignment(t *testing.T, createdBy uint) {
	t.Helper()

	start := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	end := start.Add(8 * time.Hour)

	shift := models.Shift{
		ID:        1,
		StartTime: start,
		EndTime:   end,
		CreatedBy: createdBy,
	}
	require.NoError(t, database.DB.Create(&shift).Error)
}

func TestCreateShift_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupShiftControllerTestDB(t)
	seedCompanyAndUsersForShiftTests(t)

	router := gin.New()
	router.POST("/shifts", CreateShift)

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	endTime := startTime.Add(8 * time.Hour)

	reqBody := map[string]interface{}{
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
		"created_by": uint(1),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/shifts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateShift_InvalidStartTimeFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupShiftControllerTestDB(t)
	seedCompanyAndUsersForShiftTests(t)

	router := gin.New()
	router.POST("/shifts", CreateShift)

	reqBody := map[string]interface{}{
		"start_time": "invalid-date",
		"end_time":   "2026-02-15T17:00:00Z",
		"created_by": uint(1),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/shifts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAssignWorkerToShift_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupShiftControllerTestDB(t)
	seedCompanyAndUsersForShiftTests(t)
	seedShiftForAssignment(t, 1)

	router := gin.New()
	router.POST("/shifts/assign", func(c *gin.Context) {
		// mimic JWT middleware context
		c.Set("employee_id", uint(3))
		c.Set("employeeID", uint(3)) // keep both if helpers differ
		c.Set("role", string(models.RoleSupervisor))
		AssignWorkerToShift(c)
	})

	reqBody := map[string]interface{}{
		"shift_id":    uint(1),
		"employee_id": uint(2),
		"assigned_by": uint(3),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/shifts/assign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAssignWorkerToShift_MissingFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupShiftControllerTestDB(t)
	seedCompanyAndUsersForShiftTests(t)

	router := gin.New()
	router.POST("/shifts/assign", AssignWorkerToShift)

	reqBody := map[string]interface{}{
		"shift_id": uint(1),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/shifts/assign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
