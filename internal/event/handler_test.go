package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDat(t *testing.T) {
	got, err := parseDay("2026-11-15")
	require.NoError(t, err)
	assert.Equal(t, time.Date(2026, 11, 15, 0, 0, 0, 0, time.UTC), got)

	invalidos := []string{"15/11/2026", "2026-13-01", "", "amanhã"}
	for _, s := range invalidos {
		_, err := parseDay(s)
		assert.Error(t, err, "deveria rejeitar %q", s)
	}

}
