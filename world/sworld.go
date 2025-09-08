package world

import (
	"fmt"

	"example.com/mud/parser"
	"example.com/mud/world/entities"
)

type World struct {
	RoomMap map[string]*entities.Room
}

func (w *World) Parse(player *Player, line string) string {
	cmd := parser.Parse(line)
	switch cmd.Kind {
	case parser.CommandMove:
		fmt.Println(player.Move(cmd.Params["direction"]))
	case parser.CommandTake:
		fmt.Printf("→ TAKE %q\n", cmd.Params["object"])
	case parser.CommandLook:
		fmt.Println(player.Look(cmd.Params["object"]))
	default:
		fmt.Println("I don't understand that.")
	}
	return ""
}

func (w *World) GetNeighboringRoom(r *entities.Room, direction string) *entities.Room {
	if roomId, ok := r.Exits[direction]; ok {
		room := w.RoomMap[roomId]
		return room
	}
	return nil
}
