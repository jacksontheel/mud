package world

import (
	"fmt"
	"strings"
	"sync"

	"example.com/mud/parser"
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type World struct {
	entityMap map[string]*entities.Entity
	mu        sync.RWMutex
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
	w.bus.Subscribe(player.currentRoom, player.entity, inbox)
	w.Bus().Publish(player.currentRoom, fmt.Sprintf("%s enters the room.", player.name), player.entity)
	return player
}

func (w *World) Bus() *Bus { return w.bus }

func (w *World) DisconnectPlayer(p *Player) {
	// TODO remove player from room
	w.bus.Unsubscribe(p.currentRoom, p.entity)
	w.Bus().Publish(p.currentRoom, fmt.Sprintf("%s leaves the room.", p.name), p.entity)
}

func (w *World) Publish(player *Player, message string) {
	exclude := player.entity
	w.bus.Publish(player.currentRoom, message, exclude)
}

func (w *World) Parse(player *Player, line string) string {
	cmd := parser.Parse(line)
	switch cmd.Kind {
	case parser.CommandMove:
		return player.Move(cmd.Params["direction"])
	case parser.CommandLook:
		return player.Look(cmd.Params["target"])
	case parser.CommandInventory:
		return player.Inventory()
	case parser.CommandAttack:
		return player.Attack(cmd.Params["target"], cmd.Params["instrument"])
	case parser.CommandKiss:
		return player.Kiss(cmd.Params["target"])
	default:
		return "I don't understand that."
	}
}

func (w *World) GetNeighboringRoom(r *components.Room, direction string) *entities.Entity {
	if roomId, ok := r.GetNeighboringRoomId(direction); ok {
		room := w.entityMap[roomId]
		return room
	}
	return nil
}

func (w *World) GetRoomDescription(r *entities.Entity, exclude *entities.Entity) string {
	var b strings.Builder

	room, ok := entities.GetComponent[*components.Room](r)
	if !ok {
		return "This is not a room"
	}
	roomIdentity, ok := entities.GetComponent[*components.Identity](r)
	if !ok {
		return "This room has no description"
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
	return b.String()
}
