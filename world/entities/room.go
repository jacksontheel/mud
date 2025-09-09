package entities

import (
	"strings"
)

type Room struct {
	Id              string
	Description     string
	Exits           map[string]string
	EntitiesByAlias map[string]*Entity
	AliasesByEntity map[*Entity][]string
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

	for _, e := range entities {
		aliases := getAliases(e)
		if len(aliases) == 0 {
			continue
		}
		// TODO NORMALIZE BEFORE PUTTING IN HERE
		aliasesByEntity[e] = aliases
		for _, alias := range aliases {
			key := normalizeAlias(alias)
			if key == "" {
				continue
			}
			entitiesByAlias[key] = e
		}
	}

	return &Room{
		Id:              id,
		Description:     strings.TrimSpace(description),
		Exits:           exits,
		EntitiesByAlias: entitiesByAlias,
		AliasesByEntity: aliasesByEntity,
	}
}

func (r *Room) GetDescription() string {
	var b strings.Builder

	roomDesc := strings.TrimSpace(r.Description)
	if roomDesc != "" {
		b.WriteString(roomDesc)
		b.WriteString("\n")
	}

	for e := range r.AliasesByEntity {
		if d := getDescription(e); d != "" {
			b.WriteString(d)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(r.GetExits())
	return b.String()
}

func (r *Room) GetExits() string {
	var b strings.Builder
	b.WriteString("Exits: [")

	for exit := range r.Exits {
		b.WriteString(exit)
		b.WriteString(", ")
	}

	result := strings.TrimSuffix(b.String(), ", ") + "]"
	return result
}

func normalizeAlias(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func getAliases(e *Entity) []string {
	if a, ok := Find[Aliased](e); ok {
		return a.Aliases()
	}
	return nil
}

func getDescription(e *Entity) string {
	if d, ok := Find[Descriptioned](e); ok {
		return strings.TrimSpace(d.Description())
	}
	return ""
}
