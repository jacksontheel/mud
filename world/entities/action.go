package entities

type ActionType int

const (
	ActionSay ActionType = iota
)

type Action interface {
	Id() ActionType
	Execute(ev *Event) (string, bool)
}
