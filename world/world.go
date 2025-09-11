package world

import (
	"fmt"

	"example.com/mud/parser"
	"example.com/mud/world/entities"
	"example.com/mud/world/loading"
)

type World struct {
	roomMap map[string]*entities.Room
}

func NewWorldFromJSONFile(fileName string) (*World, error) {
	rooms, err := loading.LoadRoomsFromFile("data/world.json")
	if err != nil {
		return nil, err
	}

	return &World{
		roomMap: rooms,
	}, nil
}

func (w *World) AddPlayer() *Player {
	return NewPlayer(w, w.roomMap["central"])
}

func (w *World) Parse(player *Player, line string) string {
	cmd := parser.Parse(line)
	switch cmd.Kind {
	case parser.CommandMove:
		fmt.Println(player.Move(cmd.Params["direction"]))
	case parser.CommandLook:
		fmt.Println(player.Look(cmd.Params["object"]))
	case parser.CommandAttack:
		fmt.Println(player.Attack(cmd.Params["object"]))
	case parser.CommandKiss:
		fmt.Println(player.Kiss(cmd.Params["object"]))
	default:
		fmt.Println("I don't understand that.")
	}
	return ""
}

func (w *World) GetNeighboringRoom(r *entities.Room, direction string) *entities.Room {
	if roomId, ok := r.Exits[direction]; ok {
		room := w.roomMap[roomId]
		return room
	}
	return nil
}
