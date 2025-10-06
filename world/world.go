package world

import (
	"fmt"
	"strings"

	"example.com/mud/models"
	"example.com/mud/parser"
	"example.com/mud/utils"
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type World struct {
	entityMap map[string]*entities.Entity
	bus       *Bus
}

func NewWorld(entityMap map[string]*entities.Entity) *World {
	return &World{
		entityMap: entityMap,
		bus:       NewBus(),
	}
}

func (w *World) AddPlayer(name string, inbox chan string) *Player {
	player := NewPlayer(name, w, w.entityMap["BalkingCrow"])

	if room, ok := entities.GetComponent[*components.Room](player.currentRoom); ok {
		room.AddChild(player.entity)
	}

	w.bus.Subscribe(player.currentRoom, player.entity, inbox)
	w.Publish(player, fmt.Sprintf("%s enters the room.", player.name))

	return player
}

func (w *World) Bus() *Bus { return w.bus }

func (w *World) EntitiesById() map[string]*entities.Entity { return w.entityMap }

func (w *World) DisconnectPlayer(p *Player) {
	if room, ok := entities.GetComponent[*components.Room](p.currentRoom); ok {
		room.RemoveChild(p.entity)
	}

	w.bus.Unsubscribe(p.currentRoom, p.entity)
	w.Publish(p, fmt.Sprintf("%s leaves the room.", p.name))
}

func (w *World) Publish(player *Player, message string) {
	w.bus.Publish(player.currentRoom, message, []*entities.Entity{player.entity})
}

func (w *World) Parse(player *Player, line string) (string, error) {
	cmd := parser.Parse(line)
	if cmd == nil {
		return "What in the nine hells?", nil
	}

	switch cmd.Kind {
	case "move":
		return player.Move(cmd.Params["direction"])
	case "look":
		return player.Look(cmd.Params["target"])
	case "say":
		return player.Say(cmd.Params["message"]), nil
	case "whisper":
		return player.Whisper(cmd.Params["target"], cmd.Params["message"])
	case "inventory":
		return player.Inventory()
	case "map":
		return player.Map()
	}

	// see if it has target
	if target := cmd.Params["target"]; target != "" {
		if instrument := cmd.Params["instrument"]; instrument != "" {
			response, err := player.actUponWith(cmd.Kind, target, instrument, cmd.NoMatchMessage)
			return response, err
		} else {
			response, err := player.actUpon(cmd.Kind, target, cmd.NoMatchMessage)
			return response, err
		}
	}

	return "What the hell are you talking about?", nil
}

func (w *World) GetNeighboringRoom(r *components.Room, direction string) *entities.Entity {
	if roomId, ok := r.GetNeighboringRoomId(direction); ok {
		room := w.entityMap[roomId]
		return room
	}
	return nil
}

func (w *World) GetRoomDescription(r *entities.Entity, exclude *entities.Entity) (string, error) {
	var b strings.Builder

	room, err := entities.RequireComponent[*components.Room](r)
	if err != nil {
		return "", err
	}

	formattedTitle, err := utils.FormatText(fmt.Sprintf("{'%s' | bold | red}", r.Name), map[string]string{})
	if err != nil {
		return "", fmt.Errorf("could not format room '%s' name: %w", r.Name, err)
	}

	b.WriteString(formattedTitle)
	b.WriteString("\n")

	roomDescription := strings.TrimSpace(r.Description)
	b.WriteString(roomDescription)
	b.WriteString("\n")

	for _, e := range room.GetChildren().GetChildren() {
		if e == exclude {
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
