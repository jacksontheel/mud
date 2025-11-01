package entities

type Action interface {
	Execute(ev *Event) error
}
