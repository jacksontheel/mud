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

func init() {
	registerComponentBuilder("Room", buildRoom)
}

func processComponentPrototype(def *ComponentDef) (entities.Component, error) {
	if b, ok := componentBuilders[def.Name]; ok {
		return b(def)
	}
	return nil, fmt.Errorf("could not match component name %s", def.Name)
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
			return nil, fmt.Errorf("room: unknown field %s", f.Key)
		}
	}
	return rm, nil
}
