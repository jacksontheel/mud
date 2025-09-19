package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"example.com/mud/dsl"
	"example.com/mud/world"
)

func handleConnection(conn net.Conn, gameWorld *world.World) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	fmt.Fprint(conn, "What is your name, weary adventurer? ")

	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	inbox := make(chan string, 64)
	player := gameWorld.AddPlayer(name, inbox)

	fmt.Fprintln(conn, player.OpeningMessage())

	// start consuming incoming messages
	go handleConnectionIncoming(conn, inbox)

	// notify when outgoing messages end
	done := make(chan struct{})
	go func() {
		handleConnectionOutgoing(conn, gameWorld, player)
		close(done)
	}()

	// Wait until the outgoing loop ends
	<-done
}

func handleConnectionIncoming(conn net.Conn, inbox chan string) {
	go func() {
		for msg := range inbox {
			// Use CRLF for telnet clients
			fmt.Fprint(conn, msg+"\r\n")
		}
	}()
}

func handleConnectionOutgoing(conn net.Conn, gameWorld *world.World, player *world.Player) {
	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.ToLower(line) == "quit" {
			break
		}
		fmt.Fprintln(conn, gameWorld.Parse(player, line))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
	}

	gameWorld.DisconnectPlayer(player)
	fmt.Printf("Connection closed\n")
}

func main() {
	entityMap, err := dsl.LoadEntitiesFromFile("data/world.mud")
	if err != nil {
		panic(err)
	}

	gameWorld := world.NewWorld(entityMap)

	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("MUD server listening on port 4000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, gameWorld)
	}
}
