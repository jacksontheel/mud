package dsl

import (
	"fmt"

	"example.com/mud/dsl/ast"
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type componentBuilder func(def *ast.ComponentDef) (entities.Component, error)

var componentBuilders = map[string]componentBuilder{}

func registerComponentBuilder(name string, b componentBuilder) {
	componentBuilders[name] = b
}

func init() {
	registerComponentBuilder("Room", buildRoom)
	registerComponentBuilder("Inventory", buildInventory)
	registerComponentBuilder("Container", buildContainer)
}

func processComponentPrototype(def *ast.ComponentDef) (entities.Component, error) {
	if b, ok := componentBuilders[def.Name]; ok {
		return b(def)
	}
	return nil, fmt.Errorf("could not match component name %s", def.Name)
}

func buildRoom(def *ast.ComponentDef) (entities.Component, error) {
	rm := components.NewRoom()
	for _, f := range def.Fields {
		switch f.Key {
		case "exits":
			m := f.Value.AsMap()
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

func buildInventory(_ *ast.ComponentDef) (entities.Component, error) {
	inventory := components.NewInventory()
	return inventory, nil
}

func buildContainer(def *ast.ComponentDef) (entities.Component, error) {
	container := components.NewContainer()
	for _, f := range def.Fields {
		switch f.Key {
		case "prefix":
			prefix := f.Value.String
			if prefix == nil {
				return nil, fmt.Errorf("container: prefix must be string")
			}
			container.GetChildren().SetPrefix(*prefix)
		case "revealed":
			revealed := f.Value.Bool
			if revealed == nil {
				return nil, fmt.Errorf("container: revealed must be a boolean")
			}
			container.GetChildren().SetRevealed(*revealed == "true")
		case "children":
			continue
		default:
			return nil, fmt.Errorf("room: unknown field %s", f.Key)
		}
	}
	return container, nil
}
