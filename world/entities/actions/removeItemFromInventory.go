package actions

import (
	"fmt"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type RemoveItemFromInventory struct {
	Value          string                  `json:"value"`
	InventoryOwner entities.EntitySelector `json:"inventoryOwner"`
	Item           entities.EntitySelector `json:"item"`
}

var _ entities.Action = &RemoveItemFromInventory{}

func (a *RemoveItemFromInventory) Id() string {
	return "removeItemFromInventory"
}

func (a *RemoveItemFromInventory) Execute(ev *entities.Event) (string, bool) {
	invOwner := resolve(a.InventoryOwner, ev)
	item := resolve(a.Item, ev)
	if invOwner == nil || item == nil {
		return "", false
	}

	if inv, ok := entities.GetComponent[*components.Inventory](invOwner); ok {
		inv.RemoveItem(item)
		return "", true
	}
	return "", false
}

func resolve(sel entities.EntitySelector, ev *entities.Event) *entities.Entity {
	switch sel.Type {
	case "source":
		return ev.Source
	case "instrument":
		return ev.Instrument
	case "target":
		return ev.Target
	}
	fmt.Println("This shouldn't happen.")
	return nil
}
