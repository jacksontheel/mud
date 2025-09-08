package entities

import (
	"strings"
)

type Room struct {
	Id              string
	Description     string
	Exits           map[string]string
	EntitiesByAlias map[string]Entity   // multiple entries per entities, based on number of aliases
	AliasesByEntity map[Entity][]string // get all the aliases from the entity
}

func NewRoom(id, description string, exits map[string]string, entities []Entity) *Room {
	// TODO error if nil parameters

	entitiesByAlias := make(map[string]Entity, len(entities)*2)
	aliasesByEntity := make(map[Entity][]string, len(entities))
	for _, e := range entities {
		aliases := e.GetAliases()

		aliasesByEntity[e] = aliases

		for _, alias := range aliases {
			entitiesByAlias[alias] = e
		}
	}

	return &Room{
		Id:              id,
		Description:     description,
		Exits:           exits,
		EntitiesByAlias: entitiesByAlias,
		AliasesByEntity: aliasesByEntity,
	}
}

func (r *Room) GetDescription() string {
	var b strings.Builder
	b.WriteString(r.Description)
	b.WriteString("\n")

	for item := range r.AliasesByEntity {
		b.WriteString(item.GetDescription())
		b.WriteString("\n")
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
