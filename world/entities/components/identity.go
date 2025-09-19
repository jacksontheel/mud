package components

import (
	"example.com/mud/world/entities"
)

type Identity struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Aliases     []string `json:"aliases"`
	Tags        []string `json:"tags"`
}

var _ entities.Component = &Identity{}

func (i *Identity) Id() entities.ComponentType {
	return entities.ComponentIdentity
}

func (i *Identity) Copy() entities.Component {
	return &Identity{
		Name:        i.Name,
		Description: i.Description,
		Aliases:     i.Aliases,
		Tags:        i.Tags,
	}
}
