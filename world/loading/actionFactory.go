package loading

import (
	"encoding/json"

	"example.com/mud/world/entities"
)

type actionFactory func(raw json.RawMessage) (entities.Action, error)

var actionRegistry = map[string]actionFactory{}

func init() {
	RegisterActionTypeDefaultConstructor[*entities.Say]("say")
	RegisterActionTypeDefaultConstructor[*entities.RemoveItemFromInventory]("removeItemFromInventory")
}

func RegisterActionType(typ string, ctor actionFactory) {
	actionRegistry[typ] = ctor
}

func RegisterActionTypeDefaultConstructor[T entities.Action](typ string) {
	RegisterActionType(typ, func(raw json.RawMessage) (entities.Action, error) {
		comp := new(T)
		if err := json.Unmarshal(raw, comp); err != nil {
			return nil, err
		}
		return *comp, nil
	})
}
