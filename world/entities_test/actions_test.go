package entities_test

import (
	"testing"

	"example.com/mud/world/entities"
)

func TestASay(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"basic", "Hail Brother", "Hail Brother"},
		{"empty", "", ""},
		{"special characters", "!@#$%^&*(),./';", "!@#$%^&*(),./';"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a := &entities.ASay{Text: tc.input}
			if got := a.Say(); got != tc.want {
				t.Errorf("Say() = %q, want %q", got, tc.want)
			}
		})
	}
}
