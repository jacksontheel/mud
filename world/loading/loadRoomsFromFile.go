package loading

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/mud/world/entities"
)

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

	// A list of all room ids referenced in every rooms' exits, used for validation
	referencedRooms := make(map[string]struct{}, len(rawRooms))

	for _, rr := range rawRooms {
		items, err := decodeEntities(rr.RawItems)
		if err != nil {
			return nil, fmt.Errorf("room %q: %w", rr.Id, err)
		}

		for _, v := range rr.Exits {
			referencedRooms[v] = struct{}{}
		}

		result[rr.Id] = entities.NewRoom(
			rr.Id,
			rr.Description,
			rr.Exits,
			items,
		)
	}

	// validate that every referenced room in the exits, exists
	for k := range referencedRooms {
		if _, ok := result[k]; !ok {
			return nil, fmt.Errorf("referenced room id '%s' does not exist", k)
		}
	}

	return result, nil
}

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

			comp, err := ctor(compRaw)
			if err != nil {
				return nil, fmt.Errorf("decode component %q: %w", env.Type, err)
			}

			e.Add(comp)
		}

		entitiesList = append(entitiesList, e)
	}
	return entitiesList, nil
}
