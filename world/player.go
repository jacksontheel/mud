package world

import (
	"fmt"
	"strings"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/actions"
	"example.com/mud/world/entities/components"
)

type Player struct {
	world       *World
	name        string
	entity      *entities.Entity
	currentRoom *entities.Entity
}

func NewPlayer(name string, world *World, currentRoom *entities.Entity) *Player {
	playerEntity := entities.NewEntity()
	playerEntity.Add(&components.Identity{
		Name:        name,
		Description: fmt.Sprintf("%s the brave hero is here.", name),
		Aliases:     []string{strings.ToLower(name)},
		Tags:        []string{"player"},
	})

	inventory := components.NewInventory()
	inventory.GetChildren().AddChild(getEgg())
	playerEntity.Add(inventory)

	playerEntity.Add(&components.Eventful{
		Rules: []*entities.Rule{
			{
				When: &entities.When{
					Type: "attack",
					Source: &entities.EntitySelector{
						Type:  "tag",
						Value: "player",
					},
				},
				Then: []entities.Action{
					&actions.Say{
						Text: "You beat a great big indent into this other person's head",
					},
				},
			},
		},
	})

	return &Player{
		world:       world,
		name:        name,
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
		p.world.Bus().Publish(p.currentRoom, fmt.Sprintf("%s leaves the room.", p.name), p.entity)

		currentRoom.GetChildren().RemoveChild(p.entity)
		p.currentRoom = newRoom

		if room, ok := entities.GetComponent[*components.Room](p.currentRoom); ok {
			room.GetChildren().AddChild(p.entity)
		}

		p.world.Bus().Move(p.currentRoom, p.entity)
		p.world.Bus().Publish(p.currentRoom, fmt.Sprintf("%s enters the room.", p.name), p.entity)

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
		if descriptioned, ok := entities.GetComponent[*components.Identity](target); ok { // potentially a RequireComponent
			return descriptioned.Description, nil
		}
		return fmt.Sprintf("The %s before you is undescribable.", alias), nil
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

func (p *Player) Attack(targetAlias, instrumentAlias string) (string, error) {
	if instrumentAlias == "" {
		message, err := p.actUpon(
			"attack",
			targetAlias,
			fmt.Sprintf("Now is not the time to attack %s.", targetAlias),
		)
		if err != nil {
			return "", fmt.Errorf("attack for player '%s': %w", p.name, err)
		}

		return message, nil
	}

	message, err := p.actUponWith(
		"attack",
		targetAlias,
		instrumentAlias,
		fmt.Sprintf("You reconsider attacking %s with %s, it's ridiculous.", targetAlias, instrumentAlias),
	)
	if err != nil {
		return "", fmt.Errorf("attack with instrument for player '%s': %w", p.name, err)
	}

	return message, nil
}

func (p *Player) Kiss(alias string) (string, error) {
	message, err := p.actUpon(
		"kiss",
		alias,
		fmt.Sprintf("You can be romantic with %s later.", alias),
	)
	if err != nil {
		return "", fmt.Errorf("kiss for player '%s': %w", p.name, err)
	}

	return message, nil
}
func (p *Player) actUpon(action, alias, noMatchResponse string) (string, error) {
	target, err := p.getEntityByAlias(alias)
	if err != nil {
		return "", fmt.Errorf("act upon get target for player '%s': %w", p.name, err)
	}

	if target != nil {
		if eventful, ok := entities.GetComponent[*components.Eventful](target); ok {
			response, err := eventful.OnEvent(&entities.Event{
				Type:   action,
				Source: p.entity,
				Target: target,
			})

			if err != nil {
				return "", fmt.Errorf("act upon on event for player '%s': %w", p.name, err)
			}

			if response != "" {
				return response, nil
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
		response, err := eventful.OnEvent(&entities.Event{
			Type:       action,
			Source:     p.entity,
			Instrument: instrument,
			Target:     target,
		})

		if err != nil {
			return "", fmt.Errorf("act upon with on event for player '%s': %w", p.name, err)
		}

		if response != "" {
			return response, nil
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
