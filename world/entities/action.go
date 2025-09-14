package entities

type Action interface {
	Id() string
	Execute(ev *Event) (string, bool)
}
