package models

type PatToken struct {
	Literal    string
	SlotName   string
	SlotIsRest bool
}

type Pattern struct {
	Kind           string
	Tokens         []PatToken
	NoMatchMessage string
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
