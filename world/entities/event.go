package entities

// allows us to use the bus without tightly coupling a
// specific publisher to our world model
type Publisher interface {
	Publish(room *Entity, text string, exclude []*Entity)
	PublishTo(room *Entity, recipient *Entity, text string)
}

type Event struct {
	Type       string
	Publisher  Publisher
	Room       *Entity
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
