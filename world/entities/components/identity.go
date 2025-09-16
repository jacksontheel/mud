package components

import (
	"strings"

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

func GetName(e *entities.Entity) string {
	if identity, ok := entities.GetComponent[*Identity](e); ok {
		return strings.TrimSpace(identity.Name)
	}
	return ""
}

func GetDescription(e *entities.Entity) string {
	if identity, ok := entities.GetComponent[*Identity](e); ok {
		return strings.TrimSpace(identity.Description)
	}
	return ""
}

func GetAliases(e *entities.Entity) []string {
	if identity, ok := entities.GetComponent[*Identity](e); ok {
		return identity.Aliases
	}
	return []string{}
}

func GetTags(e *entities.Entity) []string {
	if identity, ok := entities.GetComponent[*Identity](e); ok {
		return identity.Tags
	}
	return []string{}
}
