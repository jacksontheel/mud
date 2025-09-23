package commands

type Command struct {
	Kind   string
	Params map[string]string
}

type PatToken struct {
	Literal    string
	SlotName   string
	SlotIsRest bool
	SlotType   string
}

type Pattern struct {
	Kind   string
	Tokens []PatToken
}

func Lit(word string) PatToken {
	return PatToken{Literal: word}
}

func Slot(slotType, name string) PatToken {
	return PatToken{SlotName: name, SlotType: slotType}
}

func SlotRest(slotType, name string) PatToken {
	return PatToken{SlotName: name, SlotType: slotType, SlotIsRest: true}
}
