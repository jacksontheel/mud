package actions

import (
	"fmt"

	"example.com/mud/world/entities"
)

type PrintTarget int

const (
	PrintTargetUnknown PrintTarget = iota
	PrintTargetSource
	PrintTargetInstrument
	PrintTargetTarget
)

const (
	PrintTargetUnknownString    = "unknown"
	PrintTargetSourceString     = "source"
	PrintTargetInstrumentString = "instrument"
	PrintTargetTargetString     = "target"
)

func StringToPrintTarget(s string) PrintTarget {
	switch s {
	case PrintTargetSourceString:
		return PrintTargetSource
	case PrintTargetInstrumentString:
		return PrintTargetInstrument
	case PrintTargetTargetString:
		return PrintTargetTarget
	default:
		return PrintTargetUnknown
	}
}

func (pt PrintTarget) String() string {
	switch pt {
	case PrintTargetSource:
		return PrintTargetSourceString
	case PrintTargetInstrument:
		return PrintTargetSourceString
	case PrintTargetTarget:
		return PrintTargetTargetString
	default:
		return PrintTargetUnknownString
	}
}

type Print struct {
	Text   string
	Target PrintTarget
}

var _ entities.Action = &Print{}

func (p *Print) Id() entities.ActionType {
	return entities.ActionPrint
}

func (p *Print) Execute(ev *entities.Event) error {
	if ev.Publisher == nil {
		return fmt.Errorf("publisher in event may not be nil for print action")
	}

	var recipient *entities.Entity
	switch p.Target {
	case PrintTargetSource:
		recipient = ev.Source
	case PrintTargetInstrument:
		recipient = ev.Instrument
	case PrintTargetTarget:
		recipient = ev.Target
	default:
		return fmt.Errorf("invalid target '%s' for print action", p.Target.String())
	}

	ev.Publisher.PublishTo(ev.Room, recipient, p.Text)

	return nil
}
