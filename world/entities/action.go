package entities

type ActionType int

const (
	ActionPrint ActionType = iota
	ActionPublish
)

type Action interface {
	Id() ActionType
	Execute(ev *Event) error
}
