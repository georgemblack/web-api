package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/testutil"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestValidateJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config: %v", err)
	}

	ctrl := gomock.NewController(t)
	firestore := testutil.NewMockFirestoreService(ctrl)
	router := setupRouter(config, firestore)
	middleware := validateJWTMiddleware(config)

	// ==================== Test case 1: Valid token ====================
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testutil.GetJWT(config, router)))

	middleware(c)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	// ==================== Test case 2: Invalid token ====================
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer bogus")

	middleware(c)

	if w.Code != 401 {
		t.Errorf("expected status code 401, got %d", w.Code)
	}

	// ==================== Test case 3: No token ====================
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Authorization", "Bearer")

	middleware(c)

	if w.Code != 401 {
		t.Errorf("expected status code 401, got %d", w.Code)
	}

	// ==================== Test case 4: No header ====================
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	middleware(c)

	if w.Code != 401 {
		t.Errorf("expected status code 401, got %d", w.Code)
	}
}
