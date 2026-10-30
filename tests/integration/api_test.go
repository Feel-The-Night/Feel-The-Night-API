//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kisalto/Feel-The-Night-Swagger/internal/models"
)

func TestVeteranCanOnlyEditOwnGuide(t *testing.T) {
	_, adminTk := createUser(t, "admin@t.com", models.RoleAdmin)
	_, authorTk := createUser(t, "author@t.com", models.RoleVeteran)
	_, otherTk := createUser(t, "other@t.com", models.RoleVeteran)

	code, char := do(t, http.MethodPost, "/characters", adminTk, map[string]any{
		"name": "Ryu", "description": "shoto", "type": "allrounder",
	})
	require.Equal(t, http.StatusCreated, code, char)

	code, guide := do(t, http.MethodPost, "/guides", authorTk, map[string]any{
		"title": "Ryu combos", "description": "basics", "character_id": char["character_id"],
	})
	require.Equal(t, http.StatusCreated, code, guide)
	path := "/guides/" + jsonID(guide["guide_id"])

	code, _ = do(t, http.MethodPatch, path, otherTk, map[string]any{"title": "hacked"})
	assert.Equal(t, http.StatusForbidden, code, "another veteran cannot edit")

	code, _ = do(t, http.MethodPatch, path, authorTk, map[string]any{"title": "Combos v2"})
	assert.Equal(t, http.StatusOK, code, "the author can edit")

	code, _ = do(t, http.MethodPatch, path, adminTk, map[string]any{"title": "Combos v3"})
	assert.Equal(t, http.StatusOK, code, "admin can edit any guide")
}

func TestCharacterWithGuidesCannotBeDeleted(t *testing.T) {
	_, adminTk := createUser(t, "admin2@t.com", models.RoleAdmin)

	code, char := do(t, http.MethodPost, "/characters", adminTk, map[string]any{
		"name": "Ken", "description": "shoto", "type": "rushdown",
	})
	require.Equal(t, http.StatusCreated, code, char)

	code, guide := do(t, http.MethodPost, "/guides", adminTk, map[string]any{
		"title": "Ken basics", "description": "x", "character_id": char["character_id"],
	})
	require.Equal(t, http.StatusCreated, code, guide)

	code, _ = do(t, http.MethodDelete, "/characters/"+jsonID(char["character_id"]), adminTk, nil)
	assert.Equal(t, http.StatusConflict, code)
}

func TestGuideWithNonexistentCharacter(t *testing.T) {
	_, adminTk := createUser(t, "admin3@t.com", models.RoleAdmin)

	code, _ := do(t, http.MethodPost, "/guides", adminTk, map[string]any{
		"title": "Ghost", "description": "x", "character_id": 999999,
	})
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestBeginnersCannotCreateContent(t *testing.T) {
	_, tk := createUser(t, "newbie@t.com", models.RoleBeginner)

	code, _ := do(t, http.MethodPost, "/characters", tk, map[string]any{
		"name": "Chun-Li", "description": "x", "type": "zoner",
	})
	assert.Equal(t, http.StatusForbidden, code)
}

func TestLoginRateLimit(t *testing.T) {
	ip := "203.0.113.7:1234" // same IP for every attempt

	for i := 0; i < 5; i++ {
		code, _ := doFrom(t, ip, http.MethodPost, "/auth/login", "", map[string]any{
			"email": "nobody@t.com", "password": "wrong",
		})
		require.Equal(t, http.StatusUnauthorized, code, "attempt %d", i+1)
	}

	code, _ := doFrom(t, ip, http.MethodPost, "/auth/login", "", map[string]any{
		"email": "nobody@t.com", "password": "wrong",
	})
	assert.Equal(t, http.StatusTooManyRequests, code, "the 6th attempt must be blocked")
}
