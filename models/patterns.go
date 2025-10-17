package models

import (
	"fmt"
	"strings"
)

type PatToken struct {
	Literal    string
	SlotName   string
	SlotIsRest bool
}

func (pt *PatToken) String() string {
	if pt.Literal != "" {
		return pt.Literal
	}

	return fmt.Sprintf("{%s}", pt.SlotName)
}

type Pattern struct {
	Kind           string
	Tokens         []PatToken
	HelpMessage    string
	NoMatchMessage string
}

func (p *Pattern) String() string {
	var b strings.Builder

	for _, t := range p.Tokens {
		b.WriteString(t.String())
		b.WriteString(" ")
	}

	return strings.TrimSuffix(b.String(), " ")
}

func Lit(word string) PatToken {
	return PatToken{Literal: word}
}

func Slot(name string) PatToken {
	return PatToken{SlotName: name}
}

func SlotRest(name string) PatToken {
	return PatToken{SlotName: name, SlotIsRest: true}
}
