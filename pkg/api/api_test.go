package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/testutil"
	"github.com/gin-gonic/gin"
)

func TestHello(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	router := setupRouter(config)

	// Execute request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testutil.GetJWT(config, router)))
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}
}

func TestValidAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	router := setupRouter(config)
	w := httptest.NewRecorder()

	token := config.Base64UserPass()
	req, _ := http.NewRequest("POST", "/auth", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}
}

func TestValidAuthInvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	router := setupRouter(config)
	w := httptest.NewRecorder()

	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:bogus", config.APIUsername)))
	req, _ := http.NewRequest("POST", "/auth", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected status code 401, got %d", w.Code)
	}
}

func TestInvalidAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	router := setupRouter(config)
	w := httptest.NewRecorder()

	token := "Invalid"
	req, _ := http.NewRequest("POST", "/auth", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected status code 401, got %d", w.Code)
	}
}
