package entities

type ActionType int

const (
	ActionSay ActionType = iota
)

const (
	ActionWhisper ActionType = iota
)

type Action interface {
	Id() ActionType
	Execute(ev *Event) (string, bool)
}
