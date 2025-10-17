package entities

// allows us to use the bus without tightly coupling a
// specific publisher to our world model
type Publisher interface {
	Publish(room *Entity, text string, exclude []*Entity)
	PublishTo(room *Entity, recipient *Entity, text string)
}

type Event struct {
	Type         string
	Publisher    Publisher
	EntitiesById map[string]*Entity
	Room         *Entity
	Source       *Entity
	Instrument   *Entity
	Target       *Entity
	Message      string
}

type Rule struct {
	When []Condition
	Then []Action
}
