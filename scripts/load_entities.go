package scripts

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
	"github.com/dop251/goja"
)

type EntityJS struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Aliases     []string      `json:"aliases"`
	Tags        []string      `json:"tags,omitempty"`
	Components  *ComponentsJS `json:"components,omitempty"`
}

type ComponentsJS struct {
	Room *RoomComponentJS `json:"room,omitempty"`
}

type RoomComponentJS struct {
	Color    string            `json:"color"`
	Icon     string            `json:"icon"`
	Children []EntityJS        `json:"children"`
	Exits    map[string]string `json:"exits"`
}

func LoadEntitiesFromScript(path string) (map[string]*entities.Entity, error) {
	entityMap := make(map[string]*entities.Entity)

	rt := goja.New()

	rt.Set("Orbis", map[string]any{
		"load": func(v goja.Value) {
			var generic any
			_ = rt.ExportTo(v, &generic)
			raw, _ := json.Marshal(generic)

			var defs map[string]EntityJS
			_ = json.Unmarshal(raw, &defs)

			for id, def := range defs {
				entity, err := def.build()
				if err != nil {
					return
				}
				entityMap[id] = entity
			}
		},
	})

	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not get file: %w", err)
	}

	if _, err := rt.RunScript("entry", string(src)); err != nil {
		return nil, fmt.Errorf("could not run bundle: %w", err)
	}

	return entityMap, nil
}

func (ed EntityJS) build() (*entities.Entity, error) {
	entity := entities.NewEntity(ed.Name, ed.Description, ed.Aliases, ed.Tags, nil, nil)

	components, err := ed.Components.build()
	if err != nil {
		return nil, fmt.Errorf("could not build components: %w", err)
	}

	for _, c := range components {
		entity.Add(c)
	}

	return entity, nil
}

func (cjs *ComponentsJS) build() ([]entities.Component, error) {
	components := []entities.Component{}

	if cjs == nil {
		return components, nil
	}

	if cjs.Room != nil {
		room, err := cjs.Room.build()
		if err != nil {
			return nil, fmt.Errorf("could not build room: %w", err)
		}
		components = append(components, room)
	}

	return components, nil
}

func (rjs *RoomComponentJS) build() (*components.Room, error) {
	room := components.NewRoom(rjs.Icon, rjs.Color, rjs.Exits)

	for _, childJS := range rjs.Children {
		child, err := childJS.build()
		if err != nil {
			return nil, fmt.Errorf("could not build child: %w", err)
		}

		err = room.GetChildren().AddChild(child)
		if err != nil {
			return nil, fmt.Errorf("could not add child: %w", err)
		}
	}

	return room, nil
}
