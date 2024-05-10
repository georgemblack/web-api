package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/gin-gonic/gin"
)

func GetJWT(config conf.Config, router *gin.Engine) string {
	// Build and execute request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.Base64UserPass()))
	router.ServeHTTP(w, req)

	// Read token from body
	var resp types.AuthResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.Token
}
