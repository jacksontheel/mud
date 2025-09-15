package dsl

import (
	"fmt"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type componentBuilder func(def *ComponentDef) (entities.Component, error)

var componentBuilders = map[string]componentBuilder{}

func registerComponentBuilder(name string, b componentBuilder) {
	componentBuilders[name] = b
}

func processComponentNoChildren(def *ComponentDef) (entities.Component, error) {
	if b, ok := componentBuilders[def.Name]; ok {
		return b(def)
	}
	return nil, fmt.Errorf("could not match component name %s", def.Name)
}

func init() {
	registerComponentBuilder("Identity", buildIdentity)
	registerComponentBuilder("Room", buildRoom)
}

func buildIdentity(def *ComponentDef) (entities.Component, error) {
	id := &components.Identity{}
	for _, f := range def.Fields {
		switch f.Key {
		case "name":
			if f.Value.String == nil {
				return nil, fmt.Errorf("identity.name must be a string")
			}
			id.Name = *f.Value.String
		case "description":
			if f.Value.String == nil {
				return nil, fmt.Errorf("identity.description must be a string")
			}
			id.Description = *f.Value.String
		case "aliases":
			id.Aliases = f.Value.asStrings()
		case "tags":
			id.Tags = f.Value.asStrings()
		default:
			return nil, fmt.Errorf("identity: unknown field %q", f.Key)
		}
	}
	return id, nil
}

func buildRoom(def *ComponentDef) (entities.Component, error) {
	rm := components.NewRoom()
	for _, f := range def.Fields {
		switch f.Key {
		case "exits":
			m := f.Value.asMap()
			if m == nil {
				m = map[string]string{}
			}
			rm.SetExits(m)
		case "children":
			continue
		default:
			return nil, fmt.Errorf("room: unknown field %q", f.Key)
		}
	}
	return rm, nil
}
