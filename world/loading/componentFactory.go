package loading

import (
	"encoding/json"
	"fmt"

	"example.com/mud/world/entities"
)

type componentFactory func(raw json.RawMessage) (any, error)

var componentRegistry = map[string]componentFactory{}

func init() {
	RegisterComponentTypeDefaultConstructor[entities.CAliased]("aliased")
	RegisterComponentTypeDefaultConstructor[entities.CDescriptioned]("descriptioned")
	RegisterComponentType("eventful", func(raw json.RawMessage) (any, error) {
		var env rawEventful
		if err := json.Unmarshal(raw, &env); err != nil {
			return nil, fmt.Errorf("eventful envelope: %w", err)
		}

		rules := make([]entities.Rule, len(env.Rules))

		for i, r := range env.Rules {
			actions := make([]entities.Action, len(r.Then))

			for i, raw := range r.Then {
				var env rawAction
				if err := json.Unmarshal(raw, &env); err != nil {
					return nil, fmt.Errorf("rule envelope: %w", err)
				}

				if env.Type == "" {
					return nil, fmt.Errorf("action missing 'type' field")
				}
				ctor, ok := actionRegistry[env.Type]
				if !ok {
					return nil, fmt.Errorf("unknown action type %q", env.Type)
				}

				constructedAction, err := ctor(raw)
				if err != nil {
					return nil, fmt.Errorf("decode action %q: %w", env.Type, err)
				}

				actions[i] = constructedAction
			}

			rules[i] = entities.Rule{
				When: r.When,
				Then: actions,
			}
		}

		return &entities.CEventful{
			Rules: rules,
		}, nil
	})
}

func RegisterComponentType(typ string, ctor componentFactory) {
	componentRegistry[typ] = ctor
}

func RegisterComponentTypeDefaultConstructor[T any](typ string) {
	RegisterComponentType(typ, func(raw json.RawMessage) (any, error) {
		comp := new(T)
		if err := json.Unmarshal(raw, comp); err != nil {
			return nil, err
		}
		return comp, nil
	})
}
