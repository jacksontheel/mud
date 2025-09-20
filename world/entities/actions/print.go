package actions

import (
	"fmt"
	"strings"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
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

	message, err := formatText(p.Text, ev)
	if err != nil {
		return err
	}

	ev.Publisher.PublishTo(ev.Room, recipient, message)

	return nil
}

// eventually this could be its own package, handling text colors, etc
// for now it just does some simple replacements.
func formatText(s string, ev *entities.Event) (string, error) {
	if ev.Source != nil {
		identity, err := entities.RequireComponent[*components.Identity](ev.Source)
		if err != nil {
			return "", fmt.Errorf("could not format {source} for event: %w", err)
		}

		s = strings.ReplaceAll(s, "{source}", identity.Name)
	}

	if ev.Instrument != nil {
		identity, err := entities.RequireComponent[*components.Identity](ev.Instrument)
		if err != nil {
			return "", fmt.Errorf("could not format {instrument} for event: %w", err)
		}

		s = strings.ReplaceAll(s, "{instrument}", identity.Name)
	}

	if ev.Target != nil {
		identity, err := entities.RequireComponent[*components.Identity](ev.Target)
		if err != nil {
			return "", fmt.Errorf("could not format {target} for event: %w", err)
		}

		s = strings.ReplaceAll(s, "{target}", identity.Name)
	}

	return s, nil
}
