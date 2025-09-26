package entities

import "fmt"

type EventRole int

const (
	EventRoleUnknown EventRole = iota
	EventRoleSource
	EventRoleInstrument
	EventRoleTarget
	EventRoleRoom
)

const (
	EventRoleUnknownString    = "unknown"
	EventRoleSourceString     = "source"
	EventRoleInstrumentString = "instrument"
	EventRoleTargetString     = "target"
	EventRoleRoomString       = "room"
)

func ParseEventRole(s string) (EventRole, error) {
	switch s {
	case EventRoleSourceString:
		return EventRoleSource, nil
	case EventRoleInstrumentString:
		return EventRoleInstrument, nil
	case EventRoleTargetString:
		return EventRoleTarget, nil
	case EventRoleRoomString:
		return EventRoleRoom, nil
	default:
		return EventRoleUnknown, fmt.Errorf("unknown event role '%s'", s)
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
