package world

import (
	"fmt"
	"strings"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type Player struct {
	world       *World
	name        string
	entity      *entities.Entity
	currentRoom *entities.Entity
}

func NewPlayer(name string, world *World, currentRoom *entities.Entity) *Player {
	playerEntity := world.entityMap["Player"].Copy(nil)
	playerEntity.Name = name
	playerEntity.Description = fmt.Sprintf("%s the brave hero is here.", name)
	playerEntity.Aliases = []string{strings.ToLower(name)}

	return &Player{
		world:       world,
		name:        name,
		entity:      playerEntity,
		currentRoom: currentRoom,
	}
}

func (p *Player) OpeningMessage() (string, error) {
	message, err := p.world.GetRoomDescription(p.currentRoom, p.entity)
	if err != nil {
		return "", fmt.Errorf("opening message for player '%s': %w", p.name, err)
	}

	return message, nil
}

func (p *Player) Move(direction string) (string, error) {
	currentRoom, err := entities.RequireComponent[*components.Room](p.currentRoom)
	if err != nil {
		return "", fmt.Errorf("move for player '%s': %w", p.name, err)
	}

	newRoom := p.world.GetNeighboringRoom(currentRoom, direction)
	if newRoom != nil {
		p.world.Publish(p, fmt.Sprintf("%s leaves the room.", p.name))

		currentRoom.GetChildren().RemoveChild(p.entity)
		p.currentRoom = newRoom

		if room, ok := entities.GetComponent[*components.Room](p.currentRoom); ok {
			room.GetChildren().AddChild(p.entity)
		}

		p.world.Bus().Move(p.currentRoom, p.entity)
		p.world.Publish(p, fmt.Sprintf("%s enters the room.", p.name))

		return p.world.GetRoomDescription(p.currentRoom, p.entity)
	}

	return "You can't go there.", nil
}

func (p *Player) Look(alias string) (string, error) {
	if alias == "" {
		message, err := p.world.GetRoomDescription(p.currentRoom, p.entity)
		if err != nil {
			return "", fmt.Errorf("look room for player '%s': %w", p.name, err)
		}
		return message, nil
	}

	target, err := p.getEntityByAlias(alias)
	if err != nil {
		return "", fmt.Errorf("get look target for player '%s': %w", p.name, err)
	}

	if target != nil {
		description, err := target.GetDescription()
		return description, err
	}

	return "You must be going mad. That's not here.", nil
}

func (p *Player) Inventory() (string, error) {
	if inventory, ok := entities.GetComponent[*components.Inventory](p.entity); ok {
		message, err := inventory.Print()
		if err != nil {
			return "", fmt.Errorf("inventory print for player '%s': %w", p.name, err)
		}
		return message, nil
	}
	return "You couldn't possibly carry anything at all.", nil
}

func (p *Player) Say(message string) string {
	p.world.Publish(p, fmt.Sprintf("%s says, \"%s\"", p.name, message))

	return fmt.Sprintf("You say: \"%s\"", message)
}

func (p *Player) Whisper(target string, message string) (string, error) {
	recipient, err := p.getEntityByAlias(target)

	if err != nil {
		return "", fmt.Errorf("whisper error for player '%s', to target '%v': %w", p.name, recipient, err)
	}

	if recipient == nil {
		return "you utter a whisper, but nobody hears it.", nil
	}

	whisper := fmt.Sprintf("%s whispers to you: \"%s\"", p.name, message)
	p.world.bus.PublishTo(p.currentRoom, recipient, whisper)

	return fmt.Sprintf("You whisper to %s: \"%s\"", target, message), nil
}

func (p *Player) actUpon(action, alias, noMatchResponse string) (string, error) {
	target, err := p.getEntityByAlias(alias)
	if err != nil {
		return "", fmt.Errorf("act upon get target for player '%s': %w", p.name, err)
	}

	if target != nil {
		if eventful, ok := entities.GetComponent[*components.Eventful](target); ok {
			match, err := eventful.OnEvent(&entities.Event{
				Type:         action,
				Publisher:    p.world.Bus(),
				EntitiesById: p.world.EntitiesById(),
				Room:         p.currentRoom,
				Source:       p.entity,
				Target:       target,
			})

			if err != nil {
				return "", fmt.Errorf("act upon on event for player '%s': %w", p.name, err)
			}

			if match {
				return "", nil
			}
		}
		return noMatchResponse, nil
	}

	return "You must be going mad. That's not here.", nil
}

func (p *Player) actUponWith(action, targetAlias, instrumentAlias, noMatchResponse string) (string, error) {
	target, err := p.getEntityByAlias(targetAlias)
	if err != nil {
		return "", fmt.Errorf("act upon with get target for player '%s': %w", p.name, err)
	}
	if target == nil {
		return fmt.Sprintf("There is no %s here.", targetAlias), nil
	}

	instrument, err := p.getEntityByAlias(instrumentAlias)
	if err != nil {
		return "", fmt.Errorf("act upon with get instrument for player '%s': %w", p.name, err)
	}
	if instrument == nil {
		return fmt.Sprintf("You don't have %s available.", instrumentAlias), nil
	}

	if eventful, ok := entities.GetComponent[*components.Eventful](target); ok {
		match, err := eventful.OnEvent(&entities.Event{
			Type:         action,
			Publisher:    p.world.bus,
			EntitiesById: p.world.EntitiesById(),
			Room:         p.currentRoom,
			Source:       p.entity,
			Instrument:   instrument,
			Target:       target,
		})

		if err != nil {
			return "", fmt.Errorf("act upon with on event for player '%s': %w", p.name, err)
		}

		if match {
			return "", nil
		}
	}
	return noMatchResponse, nil
}

// Get entity by first looking in player's current room, then in their inventory
func (p *Player) getEntityByAlias(alias string) (*entities.Entity, error) {
	room, err := entities.RequireComponent[*components.Room](p.currentRoom)
	if err != nil {
		return nil, fmt.Errorf("getEntityByAlias for player '%s': %w", p.name, err)
	} else {
		if e, ok := room.GetChildren().GetChildByAlias(alias); ok {
			return e, nil
		}
	}

	if inventory, ok := entities.GetComponent[*components.Inventory](p.entity); ok {
		if e, ok := inventory.GetChildren().GetChildByAlias(alias); ok {
			return e, nil
		}
	}

	return nil, nil
}
