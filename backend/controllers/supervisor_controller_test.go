package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"clockit/backend/database"
	"clockit/backend/models"
	"clockit/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	if err := database.Init(path); err != nil {
		t.Fatalf("failed to init database: %v", err)
	}

	t.Cleanup(func() {
		if database.DB != nil {
			sqlDB, err := database.DB.DB()
			if err == nil && sqlDB != nil {
				_ = sqlDB.Close()
			}
			database.DB = nil
		}
		_ = os.Remove(path)
	})

	return path
}

func TestSupervisorLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	// init auth
	if err := services.InitAuth("testsecret"); err != nil {
		t.Fatalf("InitAuth failed: %v", err)
	}

	// create a company
	comp := models.Company{Name: "Test Company"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	// register supervisor
	sup, err := services.RegisterEmployee("Alice", "alice@gmail.com", "alice1234", "addr", "123", comp.ID, 10.0, models.RoleSupervisor)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	// prepare request
	body := map[string]string{"email": sup.Email, "password": "alice1234"}
	jb, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/supervisors/login", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/api/supervisors/login", SupervisorLogin)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if _, ok := resp["token"]; !ok {
		t.Fatalf("expected token in response, got %v", resp)
	}
	emp, ok := resp["employee"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected employee object, got %v", resp["employee"])
	}
	if fmt.Sprint(emp["email"]) != sup.Email {
		t.Fatalf("expected email %s, got %v", sup.Email, emp["email"])
	}
}

func TestSupervisorLogin_InvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	if err := services.InitAuth("testsecret"); err != nil {
		t.Fatalf("InitAuth failed: %v", err)
	}

	comp := models.Company{Name: "Test Company"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	_, err := services.RegisterEmployee("Bob", "bob@gmail.com", "rightpass", "addr", "123", comp.ID, 12.0, models.RoleSupervisor)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	body := map[string]string{"email": "bob@gmail.com", "password": "wrongpass"}
	jb, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/supervisors/login", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/api/supervisors/login", SupervisorLogin)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSupervisorLogin_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	if err := services.InitAuth("testsecret"); err != nil {
		t.Fatalf("InitAuth failed: %v", err)
	}

	body := map[string]string{"email": "noone@example.com", "password": "whatever"}
	jb, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/supervisors/login", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/api/supervisors/login", SupervisorLogin)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing user, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSupervisorRegister_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	comp := models.Company{Name: "RegCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	body := map[string]interface{}{
		"name":       "NewSup",
		"email":      "newsup@example.com",
		"password":   "strongpass",
		"address":    "addr",
		"phone_no":   "000",
		"company_id": comp.ID,
		"wage":       15.0,
	}
	jb, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/supervisors/register", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/api/supervisors/register", SupervisorRegister)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestSupervisorRegister_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	// missing email
	body := map[string]interface{}{
		"name":     "NoEmail",
		"password": "strongpass",
	}
	jb, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/supervisors/register", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/api/supervisors/register", SupervisorRegister)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for validation error, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetWorkerAvailability_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	// create company and worker
	comp := models.Company{Name: "AvailCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	worker, err := services.RegisterEmployee("Worker1", "w1@example.com", "pass1234", "addr", "111", comp.ID, 8.0, models.RoleWorker)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	// create availability for today's weekday
	date := "2022-03-23" // a fixed date; weekday must match created entry's day
	// compute weekday
	parsed, _ := time.Parse("2006-01-02", date)
	wa := models.WorkerAvailability{
		WorkerID:  worker.ID,
		DayOfWeek: int(parsed.Weekday()),
		StartTime: "09:00",
		EndTime:   "17:00",
	}
	if err := database.DB.Create(&wa).Error; err != nil {
		t.Fatalf("failed to create worker availability: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/supervisors/workers/availability?date="+date+"&company_id="+fmt.Sprint(comp.ID), nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/api/supervisors/workers/availability", GetWorkerAvailability)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetWorkerAvailability_InvalidDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/api/supervisors/workers/availability?date=bad-date&company_id=1", nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/api/supervisors/workers/availability", GetWorkerAvailability)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for invalid date, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestGetShiftsByCompany_Controller_ValidationAndSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	comp := models.Company{Name: "ShiftCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	sup, err := services.RegisterEmployee("Sup", "sup@example.com", "supass", "addr", "222", comp.ID, 20.0, models.RoleSupervisor)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	// create worker and shift/assignment
	worker, err := services.RegisterEmployee("W2", "w2@example.com", "pass1234", "addr", "333", comp.ID, 9.0, models.RoleWorker)
	if err != nil {
		t.Fatalf("RegisterEmployee worker failed: %v", err)
	}

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(8 * time.Hour)
	sft := models.Shift{StartTime: start, EndTime: end, CreatedBy: sup.ID, CreatedAt: time.Now()}
	if err := database.DB.Create(&sft).Error; err != nil {
		t.Fatalf("failed to create shift: %v", err)
	}

	sa := models.ShiftAssignment{ShiftID: sft.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned}
	if err := database.DB.Create(&sa).Error; err != nil {
		t.Fatalf("failed to create shift assignment: %v", err)
	}

	// valid request
	req := httptest.NewRequest(http.MethodGet, "/api/supervisors/"+fmt.Sprint(sup.ID)+"/shifts?status=Assigned", nil)
	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/api/supervisors/:employee_id/shifts", GetShiftsByCompany)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body=%s", w.Code, w.Body.String())
	}

	// invalid status filter
	req2 := httptest.NewRequest(http.MethodGet, "/api/supervisors/"+fmt.Sprint(sup.ID)+"/shifts?status=Invalid", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for invalid status, got %d, body=%s", w2.Code, w2.Body.String())
	}

	// invalid employee id
	req3 := httptest.NewRequest(http.MethodGet, "/api/supervisors/notanumber/shifts", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for invalid employee id, got %d, body=%s", w3.Code, w3.Body.String())
	}
}

func TestGetRequestedShifts_Controller_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	// create company and supervisor
	comp := models.Company{Name: "ReqCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	sup, err := services.RegisterEmployee("SupReq", "supreq@example.com", "supass", "addr", "222", comp.ID, 20.0, models.RoleSupervisor)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	// create worker and a shift with a Requested assignment
	worker, err := services.RegisterEmployee("WReq", "wreq@example.com", "pass1234", "addr", "333", comp.ID, 9.0, models.RoleWorker)
	if err != nil {
		t.Fatalf("RegisterEmployee worker failed: %v", err)
	}

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(8 * time.Hour)
	sft := models.Shift{StartTime: start, EndTime: end, CreatedBy: sup.ID, CreatedAt: time.Now()}
	if err := database.DB.Create(&sft).Error; err != nil {
		t.Fatalf("failed to create shift: %v", err)
	}

	sa := models.ShiftAssignment{ShiftID: sft.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusRequested}
	if err := database.DB.Create(&sa).Error; err != nil {
		t.Fatalf("failed to create shift assignment: %v", err)
	}

	// Make request and inject authenticated supervisor id into context before handler
	req := httptest.NewRequest(http.MethodGet, "/api/supervisors/shifts/requested", nil)
	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/api/supervisors/shifts/requested", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, sup.ID)
		GetRequestedShifts(c)
	})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body=%s", w.Code, w.Body.String())
	}

	// decode response and check there's one shift with Requested status
	var resp []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 requested shift, got %d", len(resp))
	}
	if fmt.Sprint(resp[0]["status"]) != string(models.StatusRequested) {
		t.Fatalf("expected status %s, got %v", string(models.StatusRequested), resp[0]["status"])
	}
}

func TestGetAssignedShifts_Unauthorized(t *testing.T) {
    gin.SetMode(gin.TestMode)
    setupTestDB(t)

    r := gin.New()
    r.GET("/api/supervisors/shifts/assigned", GetAssignedShifts)

    req := httptest.NewRequest(http.MethodGet, "/api/supervisors/shifts/assigned", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetAssignedShifts_SupervisorNotFound(t *testing.T) {
    gin.SetMode(gin.TestMode)
    setupTestDB(t)

    r := gin.New()
    r.GET("/api/supervisors/shifts/assigned", func(c *gin.Context) {
        // mimic auth middleware context
        c.Set("employee_id", uint(999999))
        GetAssignedShifts(c)
    })

    req := httptest.NewRequest(http.MethodGet, "/api/supervisors/shifts/assigned", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    require.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAssignedShifts_Success_Empty(t *testing.T) {
    gin.SetMode(gin.TestMode)
    setupTestDB(t)

    company := models.Company{Name: "Assigned Shift Co"}
    require.NoError(t, database.DB.Create(&company).Error)

    supervisor := models.Employee{
        Name:      "Sup A",
        Email:     "sup.assigned@test.local",
        Password:  "hashed-or-dummy",
        Role:      models.RoleSupervisor,
        CompanyID: company.ID,
        Wage:      40,
    }
    require.NoError(t, database.DB.Create(&supervisor).Error)

    r := gin.New()
    r.GET("/api/supervisors/shifts/assigned", func(c *gin.Context) {
        c.Set("employee_id", supervisor.ID)
        GetAssignedShifts(c)
    })

    req := httptest.NewRequest(http.MethodGet, "/api/supervisors/shifts/assigned", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    require.Equal(t, http.StatusOK, w.Code)

    var resp []map[string]interface{}
    require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
    require.NotNil(t, resp) // can be empty, but valid JSON array
}

func TestGetWorkersWithOvertimeHours_Controller_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	comp := models.Company{Name: "CtrlOTCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	sup, err := services.RegisterEmployee("SupCtrl", "supctrl@example.com", "supass", "addr", "222", comp.ID, 20.0, models.RoleSupervisor)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	worker, err := services.RegisterEmployee("WCtrl", "wctrl@example.com", "pass1234", "addr", "333", comp.ID, 9.0, models.RoleWorker)
	if err != nil {
		t.Fatalf("RegisterEmployee worker failed: %v", err)
	}

	now := time.Now()
	weekday := int(now.Weekday())
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-weekday, 0, 0, 0, 0, now.Location())

	for i := 1; i <= 3; i++ {
		start := weekStart.Add(time.Duration(i) * 24 * time.Hour)
		end := start.Add(8 * time.Hour)
		sft := models.Shift{StartTime: start, EndTime: end, CreatedBy: sup.ID, CreatedAt: time.Now()}
		if err := database.DB.Create(&sft).Error; err != nil {
			t.Fatalf("failed to create shift: %v", err)
		}
		sa := models.ShiftAssignment{ShiftID: sft.ID, EmployeeID: worker.ID, AssigneeID: sup.ID, AssignedAt: time.Now(), Status: models.StatusAssigned}
		if err := database.DB.Create(&sa).Error; err != nil {
			t.Fatalf("failed to create shift assignment: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/supervisors/workers/overtime", nil)
	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/api/supervisors/workers/overtime", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, sup.ID)
		GetWorkersWithOvertimeHours(c)
	})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body=%s", w.Code, w.Body.String())
	}

	var resp []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 worker over hours, got %d, body=%s", len(resp), w.Body.String())
	}
}

func TestGetWorkersWithOvertimeHours_Controller_DBFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	comp := models.Company{Name: "DBFailCo"}
	if err := database.DB.Create(&comp).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	sup, err := services.RegisterEmployee("SupDB", "supdb@example.com", "supass", "addr", "222", comp.ID, 20.0, models.RoleSupervisor)
	if err != nil {
		t.Fatalf("RegisterEmployee failed: %v", err)
	}

	if err := database.DB.Migrator().DropTable(&models.Shift{}); err != nil {
		t.Fatalf("failed to drop shifts table: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/supervisors/workers/overtime", nil)
	w := httptest.NewRecorder()
	r := gin.New()
	r.GET("/api/supervisors/workers/overtime", func(c *gin.Context) {
		c.Set(services.ContextEmployeeID, sup.ID)
		GetWorkersWithOvertimeHours(c)
	})
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error when DB query fails, got %d, body=%s", w.Code, w.Body.String())
	}
}
