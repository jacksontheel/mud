package actions

type EventRole int

const (
	EventRoleUnknown EventRole = iota
	EventRoleSource
	EventRoleInstrument
	EventRoleTarget
)

const (
	EventRoleUnknownString    = "unknown"
	EventRoleSourceString     = "source"
	EventRoleInstrumentString = "instrument"
	EventRoleTargetString     = "target"
)

func StringToEventRole(s string) EventRole {
	switch s {
	case EventRoleSourceString:
		return EventRoleSource
	case EventRoleInstrumentString:
		return EventRoleInstrument
	case EventRoleTargetString:
		return EventRoleTarget
	default:
		return EventRoleUnknown
	}
}

func (er EventRole) String() string {
	switch er {
	case EventRoleSource:
		return EventRoleSourceString
	case EventRoleInstrument:
		return EventRoleInstrumentString
	case EventRoleTarget:
		return EventRoleTargetString
	default:
		return EventRoleUnknownString
	}
}
