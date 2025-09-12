package entities

import (
	"strings"
)

type Room struct {
	id              string
	description     string
	exits           map[string]string
	entitiesByAlias map[string]*Entity
	aliasesByEntity map[*Entity][]string
}

func NewRoom(id, description string, exits map[string]string, entities []*Entity) *Room {
	if exits == nil {
		exits = make(map[string]string)
	}
	if entities == nil {
		entities = []*Entity{}
	}

	entitiesByAlias := make(map[string]*Entity, len(entities)*2)
	aliasesByEntity := make(map[*Entity][]string, len(entities))

	room := Room{
		id:              id,
		description:     strings.TrimSpace(description),
		exits:           exits,
		entitiesByAlias: entitiesByAlias,
		aliasesByEntity: aliasesByEntity,
	}

	for _, e := range entities {
		room.AddEntity(e)
	}

	return &room
}

func (r *Room) GetDescription(requester *Entity) string {
	var b strings.Builder

	roomDesc := strings.TrimSpace(r.description)
	if roomDesc != "" {
		b.WriteString(roomDesc)
		b.WriteString("\n")
	}

	for e := range r.aliasesByEntity {
		if e == requester {
			continue
		}

		if d := e.getDescription(); d != "" {
			b.WriteString(d)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(r.GetExitText())
	return b.String()
}

func (r *Room) GetEntityByAlias(alias string) (*Entity, bool) {
	entity, ok := r.entitiesByAlias[alias]
	return entity, ok
}

func (r *Room) GetNeighboringRoomId(direction string) (string, bool) {
	id, ok := r.exits[direction]
	return id, ok
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

func (r *Room) AddEntity(e *Entity) {
	aliases := e.getAliases()
	if len(aliases) == 0 {
		return
	}
	// TODO NORMALIZE BEFORE PUTTING IN HERE
	r.aliasesByEntity[e] = aliases
	for _, alias := range aliases {
		key := normalizeAlias(alias)
		if key == "" {
			continue
		}
		r.entitiesByAlias[key] = e
	}
}

func (r *Room) RemoveEntity(e *Entity) {
	delete(r.aliasesByEntity, e)

	for _, a := range e.getAliases() {
		if r.entitiesByAlias[a] == e {
			delete(r.entitiesByAlias, a)
		}
	}
}

func normalizeAlias(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
