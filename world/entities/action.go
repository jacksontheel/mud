package entities

type ActionType int

const (
	ActionPrint ActionType = iota
	ActionPublish
	ActionCopy
	ActionMove
	ActionSetField
)

type Action interface {
	Id() ActionType
	Execute(ev *Event) error
}
