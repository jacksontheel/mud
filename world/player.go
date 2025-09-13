package world

import (
	"fmt"

	"example.com/mud/world/entities"
)

type Player struct {
	world       *World
	entity      *entities.Entity
	currentRoom *entities.Room
}

func NewPlayer(name string, world *World, currentRoom *entities.Room) *Player {
	playerEntity := entities.NewEntity()
	playerEntity.Add(&entities.Descriptioned{
		Description: fmt.Sprintf("%s the brave hero is here.", name),
	})
	playerEntity.Add(&entities.Aliased{
		Aliases: []string{name, "hero"},
	})
	playerEntity.Add(&entities.Tagged{
		Tags: []string{"player"},
	})
	playerEntity.Add(&entities.Eventful{
		Rules: []entities.Rule{
			{
				When: entities.When{
					Type: "attack",
					Source: &entities.EntitySelector{
						Type: "self",
					},
				},
				Then: []entities.Action{
					&entities.Say{
						Text: "Why are you hitting yourself?",
					},
				},
			},
		},
	})
	playerEntity.Add(entities.NewInventory([]*entities.Entity{getEgg()}))

	currentRoom.AddEntity(playerEntity)

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
	egg.Add(&entities.Named{
		Name: "Egg",
	})
	egg.Add(&entities.Descriptioned{
		Description: "A bulbous, green-speckled egg.",
	})
	egg.Add(&entities.Aliased{
		Aliases: []string{"egg"},
	})
	egg.Add(&entities.Tagged{
		Tags: []string{"egg"},
	})

	return egg
}

func (p *Player) OpeningMessage() string {
	return fmt.Sprintf("You are a hero.\n%s", p.currentRoom.GetDescription(p.entity))
}

func (p *Player) Move(direction string) string {
	newRoom := p.world.GetNeighboringRoom(p.currentRoom, direction)
	if newRoom != nil {
		p.currentRoom.RemoveEntity(p.entity)
		p.currentRoom = newRoom
		p.currentRoom.AddEntity(p.entity)
		return p.currentRoom.GetDescription(p.entity)
	}
	return "You can't go there."
}

func (p *Player) Look(alias string) string {
	if alias == "" {
		return p.currentRoom.GetDescription(p.entity)
	}

	if target, ok := p.getEntityByAlias(alias); ok {
		if descriptioned, ok := entities.GetComponent[*entities.Descriptioned](target); ok {
			return descriptioned.Description
		}
		return fmt.Sprintf("The %s before you is undescribable.", alias)
	}
	return "You must be going mad. That's not here."

}

func (p *Player) Inventory() string {
	if inventory, ok := entities.GetComponent[*entities.Inventory](p.entity); ok {
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
		if eventful, ok := entities.GetComponent[*entities.Eventful](target); ok {
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

	if eventful, ok := entities.GetComponent[*entities.Eventful](target); ok {
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
	if e, ok := p.currentRoom.GetEntityByAlias(alias); ok {
		return e, true
	}

	if inventory, ok := entities.GetComponent[*entities.Inventory](p.entity); ok {
		if e, ok := inventory.GetItemByAlias(alias); ok {
			return e, true
		}
	}

	return nil, false
}
