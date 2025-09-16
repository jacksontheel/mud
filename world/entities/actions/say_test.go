package actions

import (
	"testing"

	"example.com/mud/world/entities"
	"github.com/stretchr/testify/assert"
)

func TestId(t *testing.T) {
	action := &Say{}
	assert.Equal(t, action.Id(), entities.ActionSay)
}

func TestExecute(t *testing.T) {
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
			action := &Say{Text: tc.input}
			response, _ := action.Execute(nil)

			if got := response; got != tc.want {
				t.Errorf("Say() = %q, want %q", got, tc.want)
			}
		})
	}
}
