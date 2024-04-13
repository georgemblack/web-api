package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgemblack/web-api/pkg/conf"
)

func TestHello(t *testing.T) {
	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	router := setupRouter(config)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/hello", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}
}

func TestValidAuth(t *testing.T) {
	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	router := setupRouter(config)
	w := httptest.NewRecorder()

	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.APIUsername, config.APIPassword)))
	req, _ := http.NewRequest("POST", "/auth", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}
}

func TestValidAuthInvalidCredentials(t *testing.T) {
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
