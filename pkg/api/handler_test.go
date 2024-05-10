package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
)

func TestAuthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	handler := authHandler(config)

	// ==================== Test case 1: Valid username/password ====================
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Base64UserPass()))

	handler(c)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	// Parse body into auth response
	var resp types.AuthResponse
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Error("failed to parse auth response body")
	}
	if resp.Token == "" {
		t.Error("empty jwt token in auth response")
	}

	// ==================== Test case 2: Invalid username/password ====================
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "bogus"))

	handler(c)

	if w.Code != 401 {
		t.Errorf("expected status code 401, got %d", w.Code)
	}

	// ==================== Test case 3: No header ====================
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	handler(c)

	if w.Code != 401 {
		t.Errorf("expected status code 401, got %d", w.Code)
	}
}
