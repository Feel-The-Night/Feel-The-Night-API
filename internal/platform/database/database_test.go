package database

import "testing"

func TestEscapeLike(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"texto normal", "ryu", "ryu"},
		{"porcentagem", "100%", `100\%`},
		{"underline", "ken_masters", `ken\_masters`},
		{"barra invertida", `a\b`, `a\\b`},
		{"tudo junto", `%_\`, `\%\_\\`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeLike(tt.in); got != tt.want {
				t.Errorf("EscapeLike(%q) = %q, quero %q", tt.in, got, tt.want)
			}
		})
	}
}
