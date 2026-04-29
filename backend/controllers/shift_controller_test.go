package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"clockit/backend/database"
	"clockit/backend/models"

	"clockit/backend/services"

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
	router.POST("/shifts", func(c *gin.Context) {
		// mimic JWT middleware: authenticated supervisor with ID 1
		c.Set(services.ContextEmployeeID, uint(1))
		c.Set(services.ContextEmployeeRole, string(models.RoleSupervisor))
		CreateShift(c)
	})

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
	router.POST("/shifts", func(c *gin.Context) {
		// mimic JWT middleware: authenticated supervisor with ID 1
		c.Set(services.ContextEmployeeID, uint(1))
		c.Set(services.ContextEmployeeRole, string(models.RoleSupervisor))
		CreateShift(c)
	})

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

func TestCreateShift_InvalidEndTimeFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupShiftControllerTestDB(t)
	seedCompanyAndUsersForShiftTests(t)

	router := gin.New()
	router.POST("/shifts", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, uint(1))
		CreateShift(c)
	})

	reqBody := map[string]interface{}{
		"start_time": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		"end_time":   "bad-date",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/shifts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateShift_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupShiftControllerTestDB(t)

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	body, _ := json.Marshal(map[string]string{
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   startTime.Add(8 * time.Hour).Format(time.RFC3339),
	})

	router := gin.New()
	router.POST("/shifts", CreateShift)
	req := httptest.NewRequest(http.MethodPost, "/shifts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateShift_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupShiftControllerTestDB(t)
	seedCompanyAndUsersForShiftTests(t)

	startTime := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	body, _ := json.Marshal(map[string]string{
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   startTime.Add(-1 * time.Hour).Format(time.RFC3339),
	})

	router := gin.New()
	router.POST("/shifts", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, uint(1))
		CreateShift(c)
	})
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

func TestAssignWorkerToShift_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupShiftControllerTestDB(t)

	router := gin.New()
	router.POST("/shifts/assign", AssignWorkerToShift)

	reqBody := map[string]interface{}{
		"shift_id":    uint(1),
		"employee_id": uint(2),
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/shifts/assign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestReleaseShiftForWorker_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	company := models.Company{Name: "ReleaseShiftCo"}
	require.NoError(t, database.DB.Create(&company).Error)

	sup := models.Employee{ID: 1, Name: "Sup", Email: "sup@rs.test", Role: models.RoleSupervisor, CompanyID: company.ID}
	require.NoError(t, database.DB.Create(&sup).Error)

	worker := models.Employee{ID: 2, Name: "Worker", Email: "worker@rs.test", Role: models.RoleWorker, CompanyID: company.ID}
	require.NoError(t, database.DB.Create(&worker).Error)

	start := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	end := start.Add(8 * time.Hour)
	shift := models.Shift{ID: 1, StartTime: start, EndTime: end, CreatedBy: sup.ID}
	require.NoError(t, database.DB.Create(&shift).Error)

	sa := models.ShiftAssignment{ShiftID: shift.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned}
	require.NoError(t, database.DB.Create(&sa).Error)

	router := gin.New()
	router.POST("/api/shifts/release", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, uint(worker.ID))
		c.Set(services.ContextEmployeeRole, string(models.RoleWorker))
		ReleaseShiftForWorker(c)
	})

	body := map[string]uint{"shift_id": shift.ID}
	jb, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/shifts/release", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var out models.ShiftAssignment
	require.NoError(t, database.DB.First(&out, sa.ID).Error)
	require.Equal(t, models.StatusReleased, out.Status)
}

func TestReleaseShiftForWorker_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	company := models.Company{Name: "ReleaseShiftCo"}
	require.NoError(t, database.DB.Create(&company).Error)

	worker := models.Employee{ID: 2, Name: "Worker", Email: "worker@rs.test", Role: models.RoleWorker, CompanyID: company.ID}
	require.NoError(t, database.DB.Create(&worker).Error)

	router := gin.New()
	router.POST("/api/shifts/release", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, uint(worker.ID))
		c.Set(services.ContextEmployeeRole, string(models.RoleWorker))
		ReleaseShiftForWorker(c)
	})

	body := map[string]uint{"shift_id": 9999}
	jb, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/shifts/release", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestReleaseShiftForWorker_InvalidStatus_Controller(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	company := models.Company{Name: "ReleaseShiftCo2"}
	require.NoError(t, database.DB.Create(&company).Error)

	sup := models.Employee{ID: 1, Name: "Sup2", Email: "sup2@rs.test", Role: models.RoleSupervisor, CompanyID: company.ID}
	require.NoError(t, database.DB.Create(&sup).Error)

	worker := models.Employee{ID: 2, Name: "Worker2", Email: "worker2@rs.test", Role: models.RoleWorker, CompanyID: company.ID}
	require.NoError(t, database.DB.Create(&worker).Error)

	start := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	end := start.Add(8 * time.Hour)
	shift := models.Shift{ID: 1, StartTime: start, EndTime: end, CreatedBy: sup.ID}
	require.NoError(t, database.DB.Create(&shift).Error)

	// create assignment in Requested state
	sa := models.ShiftAssignment{ShiftID: shift.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusRequested}
	require.NoError(t, database.DB.Create(&sa).Error)

	router := gin.New()
	router.POST("/api/shifts/release", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, uint(worker.ID))
		c.Set(services.ContextEmployeeRole, string(models.RoleWorker))
		ReleaseShiftForWorker(c)
	})

	body := map[string]uint{"shift_id": shift.ID}
	jb, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/shifts/release", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRequestShift_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	company := models.Company{Name: "Request Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	worker := models.Employee{
		Name:      "Worker Req",
		Email:     "worker.request@test.local",
		Password:  "password123",
		Role:      models.RoleWorker,
		CompanyID: company.ID,
		Wage:      22,
	}
	supervisor := models.Employee{
		Name:      "Supervisor Req",
		Email:     "supervisor.request@test.local",
		Password:  "password123",
		Role:      models.RoleSupervisor,
		CompanyID: company.ID,
		Wage:      45,
	}
	require.NoError(t, database.DB.Create(&worker).Error)
	require.NoError(t, database.DB.Create(&supervisor).Error)

	start := time.Now().UTC().Add(24 * time.Hour)
	end := start.Add(8 * time.Hour)
	shift := models.Shift{
		StartTime: start,
		EndTime:   end,
		CreatedBy: supervisor.ID,
	}
	require.NoError(t, database.DB.Create(&shift).Error)

	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID:    shift.ID,
		EmployeeID: worker.ID,
		AssigneeID: supervisor.ID,
		Status:     models.StatusReleased,
		AssignedAt: time.Now().UTC(),
	}).Error)

	router := gin.New()
	router.POST("/api/workers/shifts/:shift_id/request", func(c *gin.Context) {
		c.Set("employee_id", worker.ID)
		RequestShift(c)
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/workers/shifts/%d/request", shift.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var updated models.ShiftAssignment
	require.NoError(t, database.DB.Where("shift_id = ?", shift.ID).First(&updated).Error)
	require.Equal(t, models.StatusRequested, updated.Status)
	require.Equal(t, worker.ID, updated.EmployeeID)
}

func TestRequestShift_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	router := gin.New()
	router.POST("/api/workers/shifts/:shift_id/request", RequestShift)

	req := httptest.NewRequest(http.MethodPost, "/api/workers/shifts/1/request", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequestShift_BadShiftID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	router := gin.New()
	router.POST("/api/workers/shifts/:shift_id/request", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, uint(1))
		RequestShift(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/workers/shifts/bad/request", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRequestShift_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	company := models.Company{Name: "Request Missing Co"}
	require.NoError(t, database.DB.Create(&company).Error)
	worker := models.Employee{Name: "Worker Missing", Email: "worker.request.missing@test.local", Role: models.RoleWorker, CompanyID: company.ID, Wage: 22}
	require.NoError(t, database.DB.Create(&worker).Error)

	router := gin.New()
	router.POST("/api/workers/shifts/:shift_id/request", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, worker.ID)
		RequestShift(c)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/workers/shifts/99999/request", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestRejectShiftRequest_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	router := gin.New()
	router.PATCH("/api/supervisors/shifts/:shift_id/reject", RejectShiftRequest)

	req := httptest.NewRequest(http.MethodPatch, "/api/supervisors/shifts/1/reject", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRejectShiftRequest_BadShiftID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	router := gin.New()
	router.PATCH("/api/supervisors/shifts/:shift_id/reject", func(c *gin.Context) {
		c.Set("employee_id", uint(1))
		RejectShiftRequest(c)
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/supervisors/shifts/abc/reject", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRejectShiftRequest_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	company := models.Company{Name: "C1"}
	require.NoError(t, database.DB.Create(&company).Error)

	supervisor := models.Employee{
		Name:      "Sup",
		Email:     "sup.reject.ctrl@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: company.ID,
		Wage:      40,
	}
	require.NoError(t, database.DB.Create(&supervisor).Error)

	router := gin.New()
	router.PATCH("/api/supervisors/shifts/:shift_id/reject", func(c *gin.Context) {
		c.Set("employee_id", supervisor.ID)
		RejectShiftRequest(c)
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/supervisors/shifts/99999/reject", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestRejectShiftRequest_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	company := models.Company{Name: "C2"}
	require.NoError(t, database.DB.Create(&company).Error)

	supervisor := models.Employee{
		Name:      "Sup2",
		Email:     "sup2.reject.ctrl@test.local",
		Role:      models.RoleSupervisor,
		CompanyID: company.ID,
		Wage:      40,
	}
	worker := models.Employee{
		Name:      "Worker2",
		Email:     "worker2.reject.ctrl@test.local",
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

	require.NoError(t, database.DB.Create(&models.ShiftAssignment{
		ShiftID:    shift.ID,
		EmployeeID: worker.ID,
		AssigneeID: supervisor.ID,
		Status:     models.StatusRequested,
		AssignedAt: now,
	}).Error)

	router := gin.New()
	router.PATCH("/api/supervisors/shifts/:shift_id/reject", func(c *gin.Context) {
		c.Set("employee_id", supervisor.ID)
		RejectShiftRequest(c)
	})
	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/supervisors/shifts/%d/reject", shift.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var updated models.ShiftAssignment
	require.NoError(t, database.DB.Where("shift_id = ?", shift.ID).First(&updated).Error)
	require.Equal(t, models.StatusReleased, updated.Status)
	require.Equal(t, supervisor.ID, updated.AssigneeID)
}
