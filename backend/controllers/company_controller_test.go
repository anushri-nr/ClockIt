package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"clockit/backend/database"
	"clockit/backend/models"

	"github.com/gin-gonic/gin"
)

func TestCreateCompany_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	body := map[string]string{"name": "New Co", "address": "Addr"}
	jb, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/companies/create", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/api/companies/create", CreateCompany)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestCreateCompany_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	body := map[string]string{"address": "NoName"}
	jb, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/companies/create", bytes.NewReader(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := gin.New()
	r.POST("/api/companies/create", CreateCompany)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for missing name, got %d, body=%s", w.Code, w.Body.String())
	}
}

func TestListCompanies_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_ = setupTestDB(t)

	// create companies directly via DB
	c1 := models.Company{Name: "C1"}
	c2 := models.Company{Name: "C2"}
	if err := database.DB.Create(&c1).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}
	if err := database.DB.Create(&c2).Error; err != nil {
		t.Fatalf("failed to create company: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/companies/list", nil)
	w := httptest.NewRecorder()

	r := gin.New()
	r.GET("/api/companies/list", ListCompanies)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body=%s", w.Code, w.Body.String())
	}
}
