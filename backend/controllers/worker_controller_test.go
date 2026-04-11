package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fmt"
	"time"

	"clockit/backend/database"
	"clockit/backend/models"
	"clockit/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupWorkerControllerRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// routes under test
	r.POST("/api/workers/register", WorkerRegister)
	r.POST("/api/workers/login", WorkerLogin)
	r.GET("/api/workers/:employee_id/shifts", GetShiftsForWorker)
	r.POST("/api/workers/:employee_id/availability", CreateAvailability)

	return r
}

func TestWorkerRegister_BadRequest_InvalidPayload(t *testing.T) {
	r := setupWorkerControllerRouter()

	// missing required fields + invalid email
	body := `{"name":"A","email":"not-an-email","password":"123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workers/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWorkerLogin_BadRequest_InvalidPayload(t *testing.T) {
	r := setupWorkerControllerRouter()

	// missing password, invalid email
	body := `{"email":"bad-email"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workers/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetShiftsForWorker_BadRequest_InvalidEmployeeID(t *testing.T) {
	r := setupWorkerControllerRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/workers/abc/shifts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAvailability_BadRequest_InvalidEmployeeID(t *testing.T) {
	r := setupWorkerControllerRouter()

	body := `{"day_of_week":1,"start_time":"09:00:00","end_time":"17:00:00"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workers/abc/availability", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAvailability_BadRequest_InvalidBody(t *testing.T) {
	r := setupWorkerControllerRouter()

	// day_of_week out of allowed range [0..6]
	body := `{"day_of_week":9,"start_time":"09:00:00","end_time":"17:00:00"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workers/1/availability", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAvailability_WorkerNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	r := gin.New()
	r.POST("/api/workers/:employee_id/availability", CreateAvailability)

	body := `{"day_of_week":1,"start_time":"09:00:00","end_time":"17:00:00"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workers/999999/availability", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateAvailability_MalformedJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	r := gin.New()
	r.POST("/api/workers/:employee_id/availability", CreateAvailability)

	req := httptest.NewRequest(http.MethodPost, "/api/workers/1/availability", bytes.NewBufferString(`{"day_of_week":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetShiftsForWorker_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	r := gin.New()
	r.GET("/api/workers/:employee_id/shifts", GetShiftsForWorker)

	req := httptest.NewRequest(http.MethodGet, "/api/workers/999999/shifts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetShiftsForWorker_Success_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	company := models.Company{Name: "WorkerShiftCo"}
	require.NoError(t, database.DB.Create(&company).Error)

	worker := models.Employee{
		ID:        42,
		Name:      "Worker 42",
		Email:     "worker42@test.local",
		Role:      models.RoleWorker,
		CompanyID: company.ID,
		Wage:      20,
	}
	require.NoError(t, database.DB.Create(&worker).Error)

	r := gin.New()
	r.GET("/api/workers/:employee_id/shifts", GetShiftsForWorker)

	req := httptest.NewRequest(http.MethodGet, "/api/workers/42/shifts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var out []map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out))
	require.Len(t, out, 0)
}

func TestWorkerRegister_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t) // helper from supervisor_controller_test.go (same package)

	require.NoError(t, services.InitAuth("test-secret"))

	company := models.Company{Name: "Worker Register Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	r := gin.New()
	r.POST("/api/workers/register", WorkerRegister)

	body := map[string]interface{}{
		"name":       "Worker One",
		"email":      "worker.one@test.local",
		"password":   "password123",
		"address":    "123 Main St",
		"phone_no":   "1234567890",
		"company_id": company.ID,
		"wage":       22.5,
	}
	jb, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/workers/register", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.NotEmpty(t, resp["token"])
	require.NotNil(t, resp["employee"])
}

func TestWorkerRegister_RegisterEmployeeError_InvalidCompany(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	require.NoError(t, services.InitAuth("test-secret"))

	r := gin.New()
	r.POST("/api/workers/register", WorkerRegister)

	// company_id does not exist -> RegisterEmployee should fail -> 400 branch
	body := map[string]interface{}{
		"name":       "Worker Bad Company",
		"email":      "worker.badcompany@test.local",
		"password":   "password123",
		"address":    "123 Main St",
		"phone_no":   "1234567890",
		"company_id": 999999,
		"wage":       18.0,
	}
	jb, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/workers/register", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestWorkerLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t) // shared helper in controllers package

	require.NoError(t, services.InitAuth("test-secret"))

	company := models.Company{Name: "Worker Login Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	worker, err := services.RegisterEmployee(
		"Login Worker",
		"login.worker@test.local",
		"password123",
		"123 Main St",
		"9999999999",
		company.ID,
		20.0,
		models.RoleWorker,
	)
	require.NoError(t, err)
	require.NotZero(t, worker.ID)

	r := gin.New()
	r.POST("/api/workers/login", WorkerLogin)

	body := map[string]string{
		"email":    "login.worker@test.local",
		"password": "password123",
	}
	jb, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/workers/login", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.NotEmpty(t, resp["token"])
	require.NotNil(t, resp["employee"])
}

func TestGetReleasedShiftsForWorkerCompany_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	router := gin.New()
	router.GET("/api/workers/shifts/released", GetReleasedShiftsForWorkerCompany)

	req := httptest.NewRequest(http.MethodGet, "/api/workers/shifts/released", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetReleasedShiftsForWorkerCompany_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	company := models.Company{Name: "Test Co"}
	require.NoError(t, database.DB.Create(&company).Error)

	worker := models.Employee{
		Name:      "Worker One",
		Email:     "worker.released@test.local",
		Password:  "password123",
		Role:      models.RoleWorker,
		CompanyID: company.ID,
		Wage:      20,
	}
	supervisor := models.Employee{
		Name:      "Supervisor One",
		Email:     "supervisor.released@test.local",
		Password:  "password123",
		Role:      models.RoleSupervisor,
		CompanyID: company.ID,
		Wage:      40,
	}
	require.NoError(t, database.DB.Create(&worker).Error)
	require.NoError(t, database.DB.Create(&supervisor).Error)

	start := time.Now().UTC().Add(2 * time.Hour)
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
	router.GET("/api/workers/shifts/released", func(c *gin.Context) {
		c.Set("employee_id", worker.ID)
		GetReleasedShiftsForWorkerCompany(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/workers/shifts/released", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
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
