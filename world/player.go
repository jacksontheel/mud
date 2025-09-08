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

func (p *Player) Look(target string) string {
	if target == "" {
		return p.currentRoom.GetDescription()
	}
	if item, ok := p.currentRoom.EntitiesByAlias[target]; ok {
		return item.GetDescription()
	} else {
		return "you must be going mad. That's not here."
	}
}
