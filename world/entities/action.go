package entities

type ActionType int

const (
	ActionPrint ActionType = iota
	ActionPublish
	ActionCopy
	ActionMove
	ActionSetField
	ActionRevealChildren
)

type Action interface {
	Id() ActionType
	Execute(ev *Event) error
}
