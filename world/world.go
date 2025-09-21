package world

import (
	"fmt"
	"strings"

	"example.com/mud/parser"
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
	player := NewPlayer(name, w, w.entityMap["LivingRoom"])

	if room, ok := entities.GetComponent[*components.Room](player.currentRoom); ok {
		room.GetChildren().AddChild(player.entity)
	}

	w.bus.Subscribe(player.currentRoom, player.entity, inbox)
	w.Publish(player, fmt.Sprintf("%s enters the room.", player.name))

	return player
}

func (w *World) Bus() *Bus { return w.bus }

func (w *World) EntitiesById() map[string]*entities.Entity { return w.entityMap }

func (w *World) DisconnectPlayer(p *Player) {
	if room, ok := entities.GetComponent[*components.Room](p.currentRoom); ok {
		room.GetChildren().RemoveChild(p.entity)
	}

	w.bus.Unsubscribe(p.currentRoom, p.entity)
	w.Publish(p, fmt.Sprintf("%s leaves the room.", p.name))
}

func (w *World) Publish(player *Player, message string) {
	w.bus.Publish(player.currentRoom, message, []*entities.Entity{player.entity})
}

func (w *World) Parse(player *Player, line string) (string, error) {
	cmd := parser.Parse(line)
	switch cmd.Kind {
	case parser.CommandMove:
		message, err := player.Move(cmd.Params["direction"])
		return message, err
	case parser.CommandLook:
		message, err := player.Look(cmd.Params["target"])
		return message, err
	case parser.CommandInventory:
		message, err := player.Inventory()
		return message, err
	case parser.CommandAttack:
		message, err := player.Attack(cmd.Params["target"], cmd.Params["instrument"])
		return message, err
	case parser.CommandKiss:
		message, err := player.Kiss(cmd.Params["target"])
		return message, err
	case parser.CommandSay:
		return player.Say(cmd.Params["message"]), nil
	case parser.CommandWhisper:
		message, err := player.Whisper(cmd.Params["target"], cmd.Params["message"])
		return message, err
	default:
		return "I don't understand that.", nil
	}
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
	roomIdentity, ok := entities.GetComponent[*components.Identity](r)
	if !ok {
		return "This room has no description", nil
	}

	roomDescription := strings.TrimSpace(roomIdentity.Description)
	b.WriteString(roomDescription)
	b.WriteString("\n")

	for _, e := range room.GetChildren().GetChildren() {
		if e == exclude {
			continue
		}

		if eIdentity, ok := entities.GetComponent[*components.Identity](e); ok {
			b.WriteString(eIdentity.Description)
			b.WriteString("\n")
		}
	}

	b.WriteString(room.GetExitText())
	return b.String(), nil
}
