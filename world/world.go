package world

import (
	"fmt"
	"log"

	"example.com/mud/parser"
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
	"example.com/mud/world/player"
)

type World struct {
	entityMap    map[string]*entities.Entity
	startingRoom string
	bus          *Bus
}

func NewWorld(entityMap map[string]*entities.Entity, startingRoom string) *World {
	return &World{
		entityMap:    entityMap,
		startingRoom: startingRoom,
		bus:          NewBus(),
	}
}

func (w *World) EntitiesById() map[string]*entities.Entity { return w.entityMap }

func (w *World) AddPlayer(name string, inbox chan string) (*player.Player, error) {
	startingRoom, ok := w.entityMap[w.startingRoom]
	if !ok {
		log.Fatalf("add player: room '%s' does not exist in world.", w.startingRoom)
	}

	newPlayer, err := player.NewPlayer(name, w, startingRoom)
	if err != nil {
		return nil, fmt.Errorf("could not create player '%s': %w", name, err)
	}

	if room, ok := entities.GetComponent[*components.Room](newPlayer.CurrentRoom); ok {
		room.AddChild(newPlayer.Entity)
	}

	w.bus.Subscribe(newPlayer.CurrentRoom, newPlayer.Entity, inbox)
	w.Publish(newPlayer.CurrentRoom, fmt.Sprintf("%s enters the room.", newPlayer.Name), []*entities.Entity{newPlayer.Entity})

	return newPlayer, nil
}

func (w *World) DisconnectPlayer(p *player.Player) {
	if room, ok := entities.GetComponent[*components.Room](p.CurrentRoom); ok {
		room.RemoveChild(p.Entity)
	}

	w.bus.Unsubscribe(p.CurrentRoom, p.Entity)
	w.Publish(p.CurrentRoom, fmt.Sprintf("%s leaves the room.", p.Name), []*entities.Entity{p.Entity})
}

func (w *World) GetEntityById(id string) (*entities.Entity, bool) {
	entity, ok := w.entityMap[id]
	return entity, ok
}

func (w *World) Publish(room *entities.Entity, text string, exclude []*entities.Entity) {
	w.bus.Publish(room, text, exclude)
}

func (w *World) PublishTo(room *entities.Entity, recipient *entities.Entity, text string) {
	w.bus.PublishTo(room, recipient, text)
}

func (w *World) Parse(p *player.Player, line string) (string, error) {
	cmd := parser.Parse(line)
	if cmd == nil {
		return "What in the nine hells?", nil
	}

	switch cmd.Kind {
	case "move":
		return p.Move(cmd.Params["direction"])
	case "look":
		return p.Look(cmd.Params["target"])
	case "say":
		return p.Say(cmd.Params["message"]), nil
	case "whisper":
		return p.Whisper(cmd.Params["target"], cmd.Params["message"])
	case "inventory":
		return p.Inventory()
	case "map":
		return p.Map()
	}

	// see if it has target
	if target := cmd.Params["target"]; target != "" {
		if instrument := cmd.Params["instrument"]; instrument != "" {
			response, err := p.ActUponWith(cmd.Kind, target, instrument, cmd.NoMatchMessage)
			return response, err
		} else {
			response, err := p.ActUpon(cmd.Kind, target, cmd.NoMatchMessage)
			return response, err
		}
	}

	return "What the hell are you talking about?", nil
}

func (w *World) MovePlayer(p *player.Player, direction string) (string, error) {
	playerRoom, err := entities.RequireComponent[*components.Room](p.CurrentRoom)
	if err != nil {
		return "", fmt.Errorf("move for player '%s': %w", p.Name, err)
	}

	newRoom := w.getNeighboringRoom(playerRoom, direction)
	if newRoom != nil {
		w.Publish(p.CurrentRoom, fmt.Sprintf("%s leaves the room.", p.Name), []*entities.Entity{p.Entity})

		playerRoom.RemoveChild(p.Entity)
		p.CurrentRoom = newRoom

		if room, ok := entities.GetComponent[*components.Room](p.CurrentRoom); ok {
			room.AddChild(p.Entity)
		}

		w.bus.Move(p.CurrentRoom, p.Entity)
		w.Publish(p.CurrentRoom, fmt.Sprintf("%s enters the room.", p.Name), []*entities.Entity{p.Entity})

		return p.GetRoomDescription()
	}

	return "You can't go there.", nil
}

func (w *World) getNeighboringRoom(r *components.Room, direction string) *entities.Entity {
	if roomId, ok := r.GetNeighboringRoomId(direction); ok {
		room := w.entityMap[roomId]
		return room
	}
	return nil
}
