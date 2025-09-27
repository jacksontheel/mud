package commands

import "example.com/mud/models"

type CommandDefinition struct {
	Aliases []string
	Pattern []models.Pattern
}

// what should this be called
type Command struct {
	Kind           string
	Params         map[string]string
	NoMatchMessage string
}
