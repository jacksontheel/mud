package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example.com/mud/world"
)

func main() {
	gameWorld, err := world.NewWorldFromJSONFile("data/world.json")
	if err != nil {
		panic(err)
	}

	player := gameWorld.AddPlayer()

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
