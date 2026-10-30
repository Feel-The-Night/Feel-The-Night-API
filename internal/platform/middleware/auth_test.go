package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
	"github.com/kisalto/Feel-The-Night-Swagger/internal/platform/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAuthAndRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("JWT_SECRET", "a-test-secret-with-more-than-32-chars")

	r := gin.New()
	r.GET("/admin", RequireAuth(), RequireRole(models.RoleAdmin), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	token := func(role string) string {
		tk, err := jwt.GenerateAccessToken(1, "a@b.com", "x", role)
		require.NoError(t, err)
		return tk
	}

	tests := []struct {
		name   string
		header string
		want   int
	}{
		{"no header", "", http.StatusUnauthorized},
		{"made-up token", "Bearer abc.def.ghi", http.StatusUnauthorized},
		{"beginner", "Bearer " + token("beginner"), http.StatusForbidden},
		{"veteran", "Bearer " + token("veteran"), http.StatusForbidden},
		{"admin", "Bearer " + token("admin"), http.StatusOK},
		{"admin without Bearer prefix", token("admin"), http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.want, w.Code)
		})
	}

}
