package components

import (
	"strings"

	"example.com/mud/world/entities"
)

type Room struct {
	exits    map[string]string
	children *Children
}

var _ entities.Component = &Room{}
var _ entities.ComponentWithChildren = &Room{}

func NewRoom() *Room {
	return &Room{
		children: NewChildren(),
	}
}

func (r *Room) Id() entities.ComponentType {
	return entities.ComponentRoom
}

func (r *Room) Copy() entities.Component {
	// right now copying a room's children is not supported
	// return a room with no children
	return &Room{
		exits:    r.exits,
		children: NewChildren(),
	}
}

func (r *Room) GetChildren() entities.IChildren {
	return r.children
}

func (r *Room) SetExits(exits map[string]string) {
	r.exits = exits
}

func (r *Room) GetNeighboringRoomId(direction string) (string, bool) {
	roomId, ok := r.exits[direction]
	return roomId, ok
}

func (r *Room) GetExitText() string {
	var b strings.Builder
	b.WriteString("Exits: [")

	for exit := range r.exits {
		b.WriteString(exit)
		b.WriteString(", ")
	}

	result := strings.TrimSuffix(b.String(), ", ") + "]"
	return result
}
