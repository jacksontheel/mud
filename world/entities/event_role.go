package entities

import "fmt"

type EventRole int

const (
	EventRoleUnknown EventRole = iota
	EventRoleSource
	EventRoleInstrument
	EventRoleTarget
	EventRoleRoom
	EventRoleMessage
)

const (
	EventRoleUnknownString    = "unknown"
	EventRoleSourceString     = "source"
	EventRoleInstrumentString = "instrument"
	EventRoleTargetString     = "target"
	EventRoleRoomString       = "room"
	EventRoleMessageString    = "message"
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
	case EventRoleMessageString:
		return EventRoleMessage, nil
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
	case EventRoleRoom:
		return EventRoleRoomString
	case EventRoleMessage:
		return EventRoleMessageString
	default:
		return EventRoleUnknownString
	}
}
