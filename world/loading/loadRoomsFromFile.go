package loading

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/mud/world/entities"
)

// --- raw structures for decoding JSON ---

type rawRoom struct {
	Id          string            `json:"id"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
	RawItems    []rawItem         `json:"items"`
}

type rawItem struct {
	Id         string            `json:"id"`
	Components []json.RawMessage `json:"components"`
}

type rawComponent struct {
	Type string `json:"type"`
}

// --- registry machinery ---

// componentFactory creates a new empty component instance
type componentFactory func() any

var componentRegistry = map[string]componentFactory{}

// RegisterComponentType lets you plug in your concrete component structs.
// Example: RegisterComponentType("CAliases", func() any { return &entities.CAliases{} })
func RegisterComponentType(typ string, ctor componentFactory) {
	componentRegistry[typ] = ctor
}

// --- loading ---

func LoadRoomsFromFile(filename string) (map[string]*entities.Room, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var rawRooms []rawRoom
	if err := json.Unmarshal(data, &rawRooms); err != nil {
		return nil, err
	}

	result := make(map[string]*entities.Room, len(rawRooms))
	for _, rr := range rawRooms {
		items, err := decodeEntities(rr.RawItems)
		if err != nil {
			return nil, fmt.Errorf("room %q: %w", rr.Id, err)
		}

		result[rr.Id] = entities.NewRoom(
			rr.Id,
			rr.Description,
			rr.Exits,
			items,
		)
	}

	return result, nil
}

// decodeEntities turns raw JSON items into *entities.Entity
func decodeEntities(raws []rawItem) ([]*entities.Entity, error) {
	entitiesList := make([]*entities.Entity, 0, len(raws))
	for _, raw := range raws {
		e := entities.NewEntity()

		for _, compRaw := range raw.Components {
			var env rawComponent
			if err := json.Unmarshal(compRaw, &env); err != nil {
				return nil, fmt.Errorf("component envelope: %w", err)
			}
			if env.Type == "" {
				return nil, fmt.Errorf("component missing 'type' field")
			}

			ctor, ok := componentRegistry[env.Type]
			if !ok {
				return nil, fmt.Errorf("unknown component type %q", env.Type)
			}

			comp := ctor()
			if err := json.Unmarshal(compRaw, comp); err != nil {
				return nil, fmt.Errorf("decode component %q: %w", env.Type, err)
			}

			e.Add(comp)
		}

		entitiesList = append(entitiesList, e)
	}
	return entitiesList, nil
}
