//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
)

var reqCounter atomic.Int64

// nextIP returns a different IP for each request, so the rate limiter does not
// mix up requests from different tests.
func nextIP() string {
	n := reqCounter.Add(1)
	return fmt.Sprintf("10.%d.%d.%d:1234", (n>>16)&255, (n>>8)&255, n&255)
}

// doFrom sends a JSON request from a specific IP.
func doFrom(t *testing.T, ip, method, path, token string, body any) (int, map[string]any) {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.RemoteAddr = ip
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var out map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

// do sends a JSON request from a fresh IP (never hits the rate limit).
func do(t *testing.T, method, path, token string, body any) (int, map[string]any) {
	t.Helper()
	return doFrom(t, nextIP(), method, path, token, body)
}

// createUser signs up through the API, sets the role directly in the database and logs in.
func createUser(t *testing.T, email string, role models.Role) (id uint, token string) {
	t.Helper()
	code, out := do(t, http.MethodPost, "/users", "", map[string]any{
		"nickname": "tester", "email": email, "password": "password123",
	})
	require.Equal(t, http.StatusCreated, code, out)
	id = uint(out["user_id"].(float64))

	require.NoError(t, db.Model(&models.User{}).Where("user_id = ?", id).Update("role", role).Error)

	code, out = do(t, http.MethodPost, "/auth/login", "", map[string]any{
		"email": email, "password": "password123",
	})
	require.Equal(t, http.StatusOK, code, out)
	return id, out["access_token"].(string)
}

// jsonID converts an id decoded from JSON (float64) to a string for building URLs.
func jsonID(v any) string {
	return strconv.Itoa(int(v.(float64)))
}
