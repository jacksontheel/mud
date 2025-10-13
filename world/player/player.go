package player

import (
	"fmt"
	"regexp"
	"strings"

	"example.com/mud/models"
	"example.com/mud/utils"
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

var safeNameRegex = regexp.MustCompile(`[^a-zA-Z]+`)

type Player struct {
	Name        string
	Entity      *entities.Entity
	CurrentRoom *entities.Entity

	world World
}

type World interface {
	EntitiesById() map[string]*entities.Entity
	GetEntityById(id string) (*entities.Entity, bool)
	MovePlayer(p *Player, direction string) (string, error)

	Publish(room *entities.Entity, text string, exclude []*entities.Entity)
	PublishTo(room *entities.Entity, recipient *entities.Entity, text string)
}

func NewPlayer(name string, world World, currentRoom *entities.Entity) (*Player, error) {
	playerTemplate, ok := world.GetEntityById("Player")
	if !ok {
		return nil, fmt.Errorf("entity with ID 'Player' does not exist in world")
	}

	playerEntity := playerTemplate.Copy(nil)
	playerEntity.Name = name
	playerEntity.Description = fmt.Sprintf("%s the brave hero is here.", name)
	playerEntity.Aliases = []string{strings.ToLower(name)}

	return &Player{
		Name:        name,
		Entity:      playerEntity,
		CurrentRoom: currentRoom,
		world:       world,
	}, nil
}

func (p *Player) OpeningMessage() (string, error) {
	message, err := p.GetRoomDescription()
	if err != nil {
		return "", fmt.Errorf("opening message for player '%s': %w", p.Name, err)
	}

	return message, nil
}

func NameValidation(name string) string {
	if len(name) == 0 {
		return "Please, speak up! I didn't hear a name.\n"
	} else if len(name) > 20 {
		return "That's much too long to remember!\n"
	}

	testName := safeNameRegex.ReplaceAllString(name, "")

	if testName != name {
		return "I'm no good with numbers or spaces, and I only speak English!\n"
	}

	return ""
}

func (p *Player) GetRoomDescription() (string, error) {
	var b strings.Builder

	room, err := entities.RequireComponent[*components.Room](p.CurrentRoom)
	if err != nil {
		return "", err
	}

	formattedTitle, err := utils.FormatText(fmt.Sprintf("{'%s' | bold | red}", p.CurrentRoom.Name), map[string]string{})
	if err != nil {
		return "", fmt.Errorf("could not format room '%s' name: %w", p.CurrentRoom.Name, err)
	}

	b.WriteString(formattedTitle)
	b.WriteString("\n")

	roomDescription := strings.TrimSpace(p.CurrentRoom.Description)
	b.WriteString(roomDescription)
	b.WriteString("\n")

	for _, e := range room.GetChildren().GetChildren() {
		if e == p.Entity {
			continue
		}

		description, err := e.GetDescription()
		if err != nil {
			return "", err
		}

		b.WriteString(fmt.Sprintf("%s%s", models.Tab, description))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(room.GetExitText())
	return b.String(), nil
}

func (p *Player) Move(direction string) (string, error) {
	message, err := p.world.MovePlayer(p, direction)
	return message, err
}

func (p *Player) Look(alias string) (string, error) {
	if alias == "" {
		message, err := p.GetRoomDescription()
		if err != nil {
			return "", fmt.Errorf("look room for player '%s': %w", p.Name, err)
		}
		return message, nil
	}

	target, err := p.getEntityByAlias(alias)
	if err != nil {
		return "", fmt.Errorf("get look target for player '%s': %w", p.Name, err)
	}

	if target != nil {
		description, err := target.GetDescription()
		return description, err
	}

	return "You must be going mad. That's not here.", nil
}

func (p *Player) Inventory() (string, error) {
	if inventory, ok := entities.GetComponent[*components.Inventory](p.Entity); ok {
		message, err := inventory.Print()
		if err != nil {
			return "", fmt.Errorf("inventory print for player '%s': %w", p.Name, err)
		}
		return message, nil
	}
	return "You couldn't possibly carry anything at all.", nil
}

func (p *Player) ActMessage(action, message, noMatchResponse string) (string, error) {
	if eventful, ok := entities.GetComponent[*components.Eventful](p.Entity); ok {
		match, err := eventful.OnEvent(&entities.Event{
			Type:         action,
			Publisher:    p.world,
			EntitiesById: p.world.EntitiesById(),
			Room:         p.CurrentRoom,
			Source:       p.Entity,
			Target:       p.Entity,
			Message:      message,
		})

		if err != nil {
			return "", fmt.Errorf("act event for player '%s': %w", p.Name, err)
		}

		if match {
			return "", nil
		}
	}
	return noMatchResponse, nil

}

func (p *Player) ActUpon(action, targetAlias, noMatchResponse string) (string, error) {
	target, err := p.getEntityByAlias(targetAlias)
	if err != nil {
		return "", fmt.Errorf("act upon get target for player '%s': %w", p.Name, err)
	}

	if target != nil {
		if eventful, ok := entities.GetComponent[*components.Eventful](target); ok {
			match, err := eventful.OnEvent(&entities.Event{
				Type:         action,
				Publisher:    p.world,
				EntitiesById: p.world.EntitiesById(),
				Room:         p.CurrentRoom,
				Source:       p.Entity,
				Target:       target,
			})

			if err != nil {
				return "", fmt.Errorf("act upon on event for player '%s': %w", p.Name, err)
			}

			if match {
				return "", nil
			}
		}
		return noMatchResponse, nil
	}

	return "You must be going mad. That's not here.", nil
}

func (p *Player) ActUponMessage(action, targetAlias, message, noMatchResponse string) (string, error) {
	target, err := p.getEntityByAlias(targetAlias)
	if err != nil {
		return "", fmt.Errorf("act upon message get target for player '%s': %w", p.Name, err)
	}

	if target != nil {
		if eventful, ok := entities.GetComponent[*components.Eventful](target); ok {
			match, err := eventful.OnEvent(&entities.Event{
				Type:         action,
				Publisher:    p.world,
				EntitiesById: p.world.EntitiesById(),
				Room:         p.CurrentRoom,
				Source:       p.Entity,
				Target:       target,
				Message:      message,
			})

			if err != nil {
				return "", fmt.Errorf("act event for player '%s': %w", p.Name, err)
			}

			if match {
				return "", nil
			}
		}
		return noMatchResponse, nil
	}

	return "You must be going mad. That's not here.", nil
}

func (p *Player) ActUponWith(action, targetAlias, instrumentAlias, noMatchResponse string) (string, error) {
	target, err := p.getEntityByAlias(targetAlias)
	if err != nil {
		return "", fmt.Errorf("act upon with get target for player '%s': %w", p.Name, err)
	}
	if target == nil {
		return fmt.Sprintf("There is no %s here.", targetAlias), nil
	}

	instrument, err := p.getEntityByAlias(instrumentAlias)
	if err != nil {
		return "", fmt.Errorf("act upon with get instrument for player '%s': %w", p.Name, err)
	}
	if instrument == nil {
		return fmt.Sprintf("You don't have %s available.", instrumentAlias), nil
	}

	if eventful, ok := entities.GetComponent[*components.Eventful](target); ok {
		match, err := eventful.OnEvent(&entities.Event{
			Type:         action,
			Publisher:    p.world,
			EntitiesById: p.world.EntitiesById(),
			Room:         p.CurrentRoom,
			Source:       p.Entity,
			Instrument:   instrument,
			Target:       target,
		})

		if err != nil {
			return "", fmt.Errorf("act upon with on event for player '%s': %w", p.Name, err)
		}

		if match {
			return "", nil
		}
	}
	return noMatchResponse, nil
}

// Get entity by first looking in player's current room, then in their inventory
func (p *Player) getEntityByAlias(alias string) (*entities.Entity, error) {
	room, err := entities.RequireComponent[*components.Room](p.CurrentRoom)
	if err != nil {
		return nil, fmt.Errorf("getEntityByAlias for player '%s': %w", p.Name, err)
	} else {
		if e, ok := room.GetChildren().GetChildByAlias(alias); ok {
			return e, nil
		}
	}

	if inventory, ok := entities.GetComponent[*components.Inventory](p.Entity); ok {
		if e, ok := inventory.GetChildren().GetChildByAlias(alias); ok {
			return e, nil
		}
	}

	return nil, nil
}
