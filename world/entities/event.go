package entities

type Event struct {
	Type       string
	Source     *Entity
	Instrument *Entity
	Target     *Entity
}

type EntitySelector struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type When struct {
	Type       string          `json:"type"`
	Source     *EntitySelector `json:"source"`
	Instrument *EntitySelector `json:"instrument"`
}

type Rule struct {
	When *When
	Then []Action
}
