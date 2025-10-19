package entities

type ConditionType int

const (
	ConditionHasTag ConditionType = iota
	ConditionNot
	ConditionOr
	ConditionIsPresent
	ConditionEventRolesEqual
	ConditionHasChild
	ConditionFieldEquals
	ConditionMessageMatches
	ConditionExpressionTrue
)

type Condition interface {
	Check(ev *Event) (bool, error)
}
