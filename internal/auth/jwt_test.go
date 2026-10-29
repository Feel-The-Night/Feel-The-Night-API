package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessAndRefreshToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "um-segredo-de-teste-com-mais-de-32-chars")

	access, err := GenerateAccessToken(1, "a@b.com", "kisalto", "admin")
	require.NoError(t, err)

	refresh, err := GenerateRefreshToken(1)
	require.NoError(t, err)

	claims, err := ValidateAccessToken(access)
	require.NoError(t, err)
	assert.Equal(t, int64(1), claims.UserID)
	assert.Equal(t, "admin", claims.Role)

	_, err = ValidateAccessToken(refresh)
	assert.ErrorIs(t, err, ErrInvalidToken)

	_, err = ValidateRefreshToken(access)
	assert.ErrorIs(t, err, ErrInvalidToken)

}
