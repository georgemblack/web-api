package api

import (
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
