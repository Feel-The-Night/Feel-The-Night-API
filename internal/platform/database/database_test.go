package database

import "testing"

func TestEscapeLike(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"plain text", "ryu", "ryu"},
		{"percent sign", "100%", `100\%`},
		{"underline", "ken_masters", `ken\_masters`},
		{"backslash", `a\b`, `a\\b`},
		{"all together", `%_\`, `\%\_\\`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeLike(tt.in); got != tt.want {
				t.Errorf("EscapeLike(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
