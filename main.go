package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example.com/mud/world"
	"example.com/mud/world/entities"
	"example.com/mud/world/loading"
)

func main() {
	loading.RegisterComponentType("aliased", func() any { return &entities.CAliased{} })
	loading.RegisterComponentType("descriptioned", func() any { return &entities.CDescriptioned{} })
	loading.RegisterComponentType("eventful", func() any { return &entities.CEventful{} })

	rooms, err := loading.LoadRoomsFromFile("data/world.json")
	if err != nil {
		panic(err)
	}

	gameWorld := world.World{
		RoomMap: rooms,
	}

	centralRoom := rooms["central"]
	player := world.NewPlayer(&gameWorld, centralRoom)

	in := bufio.NewScanner(os.Stdin)
	fmt.Println(player.OpeningMessage())
	for {
		fmt.Print("> ")
		if !in.Scan() {
			break
		}
		line := strings.TrimSpace(in.Text())
		if line == "" {
			continue
		}
		if strings.ToLower(line) == "quit" {
			break
		}
		gameWorld.Parse(player, line)
	}
}
