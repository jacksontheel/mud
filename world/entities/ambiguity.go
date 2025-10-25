package entities

import "errors"

type AmbiguitySlot struct {
	Role    string
	Prompt  string
	Matches []AmbiguityOption
}

type AmbiguityOption struct {
	Text   string
	Entity *Entity
}

var ErrTargetAmbiguous = errors.New("target ambiguous")

type AmbiguityError struct {
	Slots   []AmbiguitySlot
	Execute func(map[string]*Entity) (string, error)
}

func (e AmbiguityError) Error() string  { return ErrTargetAmbiguous.Error() }
func (e *AmbiguityError) Unwrap() error { return ErrTargetAmbiguous }

type PendingAction struct {
	Ambiguity *AmbiguityError
	StepIndex int
	Selected  map[string]int
}
