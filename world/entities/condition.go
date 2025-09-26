package entities

type ConditionType int

const (
	ConditionHasTag ConditionType = iota
	ConditionNot
	ConditionIsPresent
	ConditionEventRolesEqual
)

type Condition interface {
	Check(ev *Event) (bool, error)
}
