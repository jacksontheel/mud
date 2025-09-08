package loading

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/mud/world/entities"
)

type entityFactory func() entities.Entity

var entityRegistry = map[string]entityFactory{}

func RegisterEntityType(typ string, ctor entityFactory) {
	entityRegistry[typ] = ctor
}

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
		items, err := decodeEntities(rr.Entities)
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

func decodeEntities(raws []json.RawMessage) ([]entities.Entity, error) {
	entities := make([]entities.Entity, 0, len(raws))
	for _, raw := range raws {
		// look at type
		var env entityEnvelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, fmt.Errorf("item envelope: %w", err)
		}
		if env.Type == "" {
			return nil, fmt.Errorf("item missing 'type' field")
		}

		// find constructor
		ctor, ok := entityRegistry[env.Type]
		if !ok {
			return nil, fmt.Errorf("unknown item type %q", env.Type)
		}

		// decode
		inst := ctor()
		if err := json.Unmarshal(raw, inst); err != nil {
			return nil, fmt.Errorf("decode %q: %w", env.Type, err)
		}

		entities = append(entities, inst)
	}
	return entities, nil
}
