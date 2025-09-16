package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example.com/mud/dsl"
	"example.com/mud/world"
)

func main() {
	entityMap, err := dsl.LoadEntitiesFromFile("data/world.mud")
	if err != nil {
		panic(err)
	}

	gameWorld := world.NewWorld(entityMap)

	player := gameWorld.AddPlayer("Craig")

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
