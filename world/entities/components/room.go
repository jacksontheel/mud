package components

import "example.com/mud/world/entities"

type Room struct {
	exits          map[string]string
	childByAlias   map[string]*entities.Entity
	aliasesByChild map[*entities.Entity][]string
}

var _ entities.Component = &Room{}

func NewRoom() *Room {
	return &Room{
		childByAlias:   make(map[string]*entities.Entity),
		aliasesByChild: make(map[*entities.Entity][]string),
	}
}

func (r *Room) Copy() entities.Component {
	// right now copying a room's children is not supported
	// return a room with no children
	return &Room{
		exits:          r.exits,
		childByAlias:   make(map[string]*entities.Entity),
		aliasesByChild: make(map[*entities.Entity][]string),
	}
}

func (r *Room) SetExits(exits map[string]string) {
	r.exits = exits
}

func (r *Room) AddChild(child *entities.Entity) {
	aliases := GetAliases(child)
	if len(aliases) == 0 {
		return
	}
	for _, alias := range aliases {
		r.aliasesByChild[child] = append(r.aliasesByChild[child], alias)
		r.childByAlias[alias] = child
	}
}

func (r *Room) RemoveChild(child *entities.Entity) {
	aliases, ok := r.aliasesByChild[child]
	if !ok {
		return
	}

	for _, alias := range aliases {
		delete(r.childByAlias, alias) // for each alias, delete from itemsByAlias
	}
	delete(r.aliasesByChild, child) // delete entry from aliasesByItem
}

func (r *Room) GetNeighboringRoomId(direction string) (string, bool) {
	roomId, ok := r.exits[direction]
	return roomId, ok
}

func (r *Room) GetDescription() string {

}
