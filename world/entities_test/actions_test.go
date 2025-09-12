package entities_test

import (
	"testing"

	"example.com/mud/world/entities"
)

func TestASay(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Hail Brother", "Hail Brother"},
		{"", ""},
		{"!@#$%^&*(),./';", "!@#$%^&*(),./';"},
	}

	for _, c := range cases {
		a := &entities.ASay{Text: c.input}
		if got := a.Say(); got != c.want {
			t.Errorf("Say() = %q, want %q", got, c.want)
		}
	}
}
