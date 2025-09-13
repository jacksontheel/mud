package entities

type Action interface {
	Id() string
}

type Say struct {
	Text string `json:"text"`
}

var _ Action = &Say{}

func (a *Say) Id() string {
	return "say"
}

type RemoveItemFromInventory struct {
	Value          string         `json:"value"`
	InventoryOwner EntitySelector `json:"inventoryOwner"`
	Item           EntitySelector `json:"item"`
}

var _ Action = &RemoveItemFromInventory{}

func (a *RemoveItemFromInventory) Id() string {
	return "removeItemFromInventory"
}

func (a *RemoveItemFromInventory) RemoveItemFromInventory(eInventory *Entity, eToRemove *Entity) {
	if inventory, ok := GetComponent[*Inventory](eInventory); ok {
		inventory.RemoveItem(eToRemove)
	}
}
