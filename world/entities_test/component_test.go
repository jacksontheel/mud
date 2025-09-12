package entities_test

import (
	"testing"

	"example.com/mud/world/entities"
)

/*
	component.go consists of two main structs with minimal logic
	the goal of testing these structs is to:
		- Ensure the return of an expected value
		- Ensure the struct satisfies the interface contract (already checked for at compile time but whatever)
*/

func TestComponentInterfaceImplementations(t *testing.T) {
	t.Run("Descriptioned", func(t *testing.T) {
		cases := []struct {
			name  string
			input string
			want  string
		}{
			{"basic", "A hole only big enough for a frog", "A hole only big enough for a frog"},
			{"empty", "", ""},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				// Explicit variable declaration ensures value implements the Descriptioned interface
				var d entities.Descriptioned = &entities.CDescriptioned{Text: tc.input}
				got := d.Description()
				if got != tc.want {
					t.Errorf("Description() = %q, want %q", got, tc.want)
				}
			})
		}
	})

	t.Run("Aliased", func(t *testing.T) {
		cases := []struct {
			name  string
			input []string
			want  []string
		}{
			{"basic", []string{"guy", "dude"}, []string{"guy", "dude"}},
			{"empty", []string{}, []string{}},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				// Explicit variable declaration ensures value implements the Descriptioned interface
				var a entities.Aliased = &entities.CAliased{Text: tc.input}
				got := a.Aliases()
				if len(got) != len(tc.want) {
					t.Fatalf("Aliases() length = %d, want %d", len(got), len(tc.want))
				}
				for i := range got {
					if got[i] != tc.want[i] {
						t.Errorf("Aliases()[%d] = %q, want %q", i, got[i], tc.want[i])
					}
				}
			})
		}
	})
}
