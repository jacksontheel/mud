package world

import (
	"fmt"

	"example.com/mud/parser"
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

type World struct {
	entityMap map[string]*entities.Entity
}

func NewWorld(entityMap map[string]*entities.Entity) *World {
	return &World{
		entityMap: entityMap,
	}
}

func (w *World) AddPlayer(name string) *Player {
	return NewPlayer(name, w, w.entityMap["LivingRoom"])
}

func (w *World) Parse(player *Player, line string) string {
	cmd := parser.Parse(line)
	switch cmd.Kind {
	case parser.CommandMove:
		fmt.Println(player.Move(cmd.Params["direction"]))
	case parser.CommandLook:
		fmt.Println(player.Look(cmd.Params["target"]))
	case parser.CommandInventory:
		fmt.Println(player.Inventory())
	case parser.CommandAttack:
		fmt.Println(player.Attack(cmd.Params["target"], cmd.Params["instrument"]))
	case parser.CommandKiss:
		fmt.Println(player.Kiss(cmd.Params["target"]))
	default:
		fmt.Println("I don't understand that.")
	}
	return ""
}

func (w *World) GetNeighboringRoom(r *components.Room, direction string) *entities.Entity {
	if roomId, ok := r.GetNeighboringRoomId(direction); ok {
		room := w.entityMap[roomId]
		return room
	}
	return nil
}
