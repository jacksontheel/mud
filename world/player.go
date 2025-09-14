package world

import (
	"fmt"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type Player struct {
	world       *World
	entity      *entities.Entity
	currentRoom *entities.Entity
}

func NewPlayer(name string, world *World, currentRoom *entities.Entity) *Player {
	playerEntity := entities.NewEntity()
	playerEntity.Add(&components.Identity{
		Name:        name,
		Description: fmt.Sprintf("%s the brave hero is here.", name),
		Aliases:     []string{name, "hero"},
		Tags:        []string{"player"},
	})

	playerEntity.Add(components.NewInventory([]*entities.Entity{getEgg()}))

	if room, ok := entities.GetComponent[*components.Room](playerEntity); ok {
		room.AddChild(playerEntity)
	}

	// TODO ERROR HANDLING
	return &Player{
		world:       world,
		entity:      playerEntity,
		currentRoom: currentRoom,
	}
}

// temporary function to test inventory
func getEgg() *entities.Entity {
	egg := entities.NewEntity()
	egg.Add(&components.Identity{
		Name:        "Egg",
		Description: "A bulbous, green-speckled egg.",
		Aliases:     []string{"egg"},
		Tags:        []string{"egg"},
	})

	return egg
}

func (p *Player) OpeningMessage() string {
	return "You are a hero."
}

func (p *Player) Move(direction string) string {
	currentRoom, ok := entities.GetComponent[*components.Room](p.currentRoom)
	if !ok {
		return "You cannot hope to leave this room, it isn't a room at all."
	}

	newRoom := p.world.GetNeighboringRoom(currentRoom, direction)
	if newRoom != nil {
		currentRoom.RemoveChild(p.entity)
		p.currentRoom = newRoom

		if room, ok := entities.GetComponent[*components.Room](p.currentRoom); ok {
			room.AddChild(p.entity)
		}

		if identity, ok := entities.GetComponent[*components.Identity](p.currentRoom); ok {
			return identity.Description
		}

		return "You enter a room that cannot be described."
	}

	return "You can't go there."
}

func (p *Player) Look(alias string) string {
	if alias == "" {
		if identity, ok := entities.GetComponent[*components.Identity](p.currentRoom); ok {
			return identity.Description
		}
		return "The room you are in cannot be described."
	}

	if target, ok := p.getEntityByAlias(alias); ok {
		if descriptioned, ok := entities.GetComponent[*components.Identity](target); ok {
			return descriptioned.Description
		}
		return fmt.Sprintf("The %s before you is undescribable.", alias)
	}
	return "You must be going mad. That's not here."

}

func (p *Player) Inventory() string {
	if inventory, ok := entities.GetComponent[*components.Inventory](p.entity); ok {
		return inventory.Print()
	}
	return "You couldn't possibly carry anything at all."
}

func (p *Player) Attack(targetAlias, instrumentAlias string) string {
	if instrumentAlias == "" {
		return p.actUpon(
			"attack",
			targetAlias,
			fmt.Sprintf("Now is not the time to attack %s.", targetAlias),
		)
	}

	return p.actUponWith(
		"attack",
		targetAlias,
		instrumentAlias,
		fmt.Sprintf("You reconsider attacking %s with %s, it's ridiculous.", targetAlias, instrumentAlias),
	)
}

func (p *Player) Kiss(alias string) string {
	return p.actUpon(
		"kiss",
		alias,
		fmt.Sprintf("You can be romantic with %s later.", alias),
	)
}

func (p *Player) actUpon(action, alias, noMatchResponse string) string {
	if target, ok := p.getEntityByAlias(alias); ok {
		if eventful, ok := entities.GetComponent[*components.Eventful](target); ok {
			response, ok := eventful.OnEvent(&entities.Event{
				Type:   action,
				Source: p.entity,
				Target: target,
			})

			if ok {
				return response
			}
		}
		return noMatchResponse
	}
	return "You must be going mad. That's not here."
}

func (p *Player) actUponWith(action, targetAlias, instrumentAlias, noMatchResponse string) string {
	target, targetOk := p.getEntityByAlias(targetAlias)
	if !targetOk {
		return fmt.Sprintf("There is no %s here.", targetAlias)
	}

	instrument, instrumentOk := p.getEntityByAlias(instrumentAlias)
	if !instrumentOk {
		return fmt.Sprintf("You don't have %s available.", instrumentAlias)
	}

	if eventful, ok := entities.GetComponent[*components.Eventful](target); ok {
		response, ok := eventful.OnEvent(&entities.Event{
			Type:       action,
			Source:     p.entity,
			Instrument: instrument,
			Target:     target,
		})

		if ok {
			return response
		}
	}
	return noMatchResponse
}

// Get entity by first looking in player's current room, then in their inventory
func (p *Player) getEntityByAlias(alias string) (*entities.Entity, bool) {
	// TODO broke this
	return nil, false
}
