package entities

type ConditionType int

const (
	ConditionHasTag ConditionType = iota
	ConditionNot
	ConditionIsPresent
	ConditionEventRolesEqual
	ConditionHasChild
	ConditionFieldEquals
)

type Condition interface {
	Check(ev *Event) (bool, error)
}
