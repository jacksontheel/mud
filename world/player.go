package world

import (
	"fmt"

	"example.com/mud/world/entities"
)

type Player struct {
	world       *World
	currentRoom *entities.Room
}

func NewPlayer(world *World, currentRoom *entities.Room) *Player {
	// TODO ERROR HANDLING
	return &Player{
		world:       world,
		currentRoom: currentRoom,
	}
}

func (p *Player) OpeningMessage() string {
	return fmt.Sprintf("You are a hero.\n%s", p.currentRoom.GetDescription())
}

func (p *Player) Move(direction string) string {
	newRoom := p.world.GetNeighboringRoom(p.currentRoom, direction)
	if newRoom != nil {
		p.currentRoom = newRoom
		return p.currentRoom.GetDescription()
	}
	return "You can't go there."
}

func (p *Player) Look(alias string) string {
	if alias == "" {
		return p.currentRoom.GetDescription()
	}

	if target, ok := p.currentRoom.EntitiesByAlias[alias]; ok {
		if eventful, ok := entities.Find[entities.Descriptioned](target); ok {
			return eventful.Description()
		}
		return fmt.Sprintf("The %s before you is undescribable.", alias)
	}
	return "you must be going mad. That's not here."

}

func (p *Player) Attack(alias string) string {
	return p.actUpon(
		"attack",
		alias,
		fmt.Sprintf("Now is not the time to attack %s.", alias),
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
	if target, ok := p.currentRoom.EntitiesByAlias[alias]; ok {
		if eventful, ok := entities.Find[entities.Eventful](target); ok {
			response, ok := eventful.OnEvent(&entities.Event{
				Type:   action,
				Target: target,
			})

			if ok {
				return response
			}
		}
		return noMatchResponse
	}
	return "you must be going mad. That's not here."
}
