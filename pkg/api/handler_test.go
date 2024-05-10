package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/testutil"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
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

	if w.Code != http.StatusOK {
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

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code 401, got %d", w.Code)
	}

	// ==================== Test case 3: No header ====================
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	handler(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code 401, got %d", w.Code)
	}
}

func TestGetLikesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	fs := testutil.NewMockFirestoreService(ctrl)
	handler := getLikesHandler(fs)

	// ==================== Test case 1: Valid request ====================
	fs.EXPECT().GetLikes().Return(testutil.NewLikes(), nil)

	// Execute handler
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/likes", nil)
	handler(c)

	// Check
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	// ==================== Test case 2: Valid request, internal error ====================
	fs.EXPECT().GetLikes().Return(testutil.NewLikes(), errors.New("ope"))

	// Execute handler
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/likes", nil)
	handler(c)

	// Check
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", w.Code)
	}
}

func TestGetLikeHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := gomock.NewController(t)
	fs := testutil.NewMockFirestoreService(ctrl)
	handler := getLikeHandler(fs)

	// ==================== Test case 1: Valid request ====================
	like := testutil.NewLike()
	fs.EXPECT().GetLike(like.ID).Return(like, nil)

	// Execute handler
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: like.ID})
	handler(c)

	// Check
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", w.Code)
	}

	// ==================== Test case 2: Missing like ID ====================
	like = testutil.NewLike()
	fs.EXPECT().GetLike(like.ID).Times(0)

	// Execute handler
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	handler(c)

	// Check
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code 400, got %d", w.Code)
	}
}
