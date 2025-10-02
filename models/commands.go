package models

type CommandDefinition struct {
	Name     string
	Aliases  []string
	Patterns []CommandPattern
}

type CommandPattern struct {
	Tokens         []PatToken
	NoMatchMessage string
}

type Command struct {
	Kind           string
	Params         map[string]string
	NoMatchMessage string
}
