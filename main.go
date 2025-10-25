package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"example.com/mud/config"
	"example.com/mud/dsl"
	"example.com/mud/parser/commands"
	"example.com/mud/world"
	"example.com/mud/world/entities"
	"example.com/mud/world/player"
)

func handleConnection(conn net.Conn, gameWorld *world.World) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var name string

	for {

		if _, err := fmt.Fprint(conn, "What is your name, weary adventurer? "); err != nil {
			return
		}

		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)

		vdn := player.NameValidation(name)
		if vdn != "" {
			fmt.Fprint(conn, vdn)
			continue
		}
		break
	}

	inbox := make(chan string, 64)
	player, err := gameWorld.AddPlayer(name, inbox)
	if err != nil {
		err := fmt.Errorf("error adding player: %w", err)

		fmt.Println(err.Error())
		fmt.Fprintln(conn, err.Error())
		return
	}

	message, err := player.OpeningMessage()
	if err != nil {
		err := fmt.Errorf("error printing opening message: %w", err)

		fmt.Println(err.Error())
		fmt.Fprintln(conn, err.Error())

		return
	} else {
		fmt.Fprintln(conn, message)
	}

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

func handleConnectionOutgoing(conn net.Conn, gameWorld *world.World, player *player.Player) {
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

		// check if player has pending multi-part messages
		if p := player.Pending; p != nil {
			n, nerr := strconv.Atoi(line)
			if nerr == nil {
				slot := p.Ambiguity.Slots[p.StepIndex]
				if n >= 1 && n <= len(slot.Matches) {
					p.Selected[slot.Role] = n - 1
					p.StepIndex++
					if p.StepIndex < len(p.Ambiguity.Slots) {
						// continue prompting current slot
						promptCurrentSlot(conn, p)
					} else {
						chosen := make(map[string]*entities.Entity, len(p.Selected))
						for _, s := range p.Ambiguity.Slots {
							idx := p.Selected[s.Role]
							chosen[s.Role] = s.Matches[idx].Entity
						}
						out, execErr := p.Ambiguity.Execute(chosen)
						player.Pending = nil
						if execErr != nil {
							fmt.Fprintln(conn, execErr.Error())
						} else if out != "" {
							fmt.Fprintln(conn, out)
						}
					}
				}
				continue
			}
			// non-number input falls through and removes pending action
			player.Pending = nil
		}

		message, err := gameWorld.Parse(player, line)
		if err != nil {
			var amb *entities.AmbiguityError
			if errors.As(err, &amb) {
				player.Pending = &entities.PendingAction{
					Ambiguity: amb,
					StepIndex: 0,
					Selected:  map[string]int{},
				}
				promptCurrentSlot(conn, player.Pending)
				continue
			}

			err := fmt.Errorf("error received: %w", err)

			fmt.Println(err.Error())
			fmt.Fprintln(conn, err.Error())
		} else if message != "" {
			fmt.Fprintln(conn, message)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Connection error:", err)
	}

	gameWorld.DisconnectPlayer(player)
	fmt.Printf("Connection closed\n")
}

func promptCurrentSlot(conn net.Conn, p *entities.PendingAction) {
	slot := p.Ambiguity.Slots[p.StepIndex]
	fmt.Fprintln(conn, slot.Prompt)
	for i, opt := range slot.Matches {
		fmt.Fprintf(conn, "  %d) %s\r\n", i+1, opt.Text)
	}
}

func main() {
	// load configuration file
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	entityMap, cmds, err := dsl.LoadEntitiesFromDirectory("data/")
	if err != nil {
		log.Fatalf("failed to load DSL entities: %v", err)
	}

	// validate starting room exists in entity map
	if _, ok := entityMap[cfg.StartingRoom]; !ok {
		log.Fatalf("room '%s' does not exist in world.", cfg.StartingRoom)
	}

	if err := commands.RegisterBuiltInCommands(); err != nil {
		log.Fatalf("failed to register built-in commands: %v", err)
	}

	if err := commands.RegisterCommands(cmds); err != nil {
		log.Fatalf("failed to register DSL commands: %v", err)
	}

	gameWorld := world.NewWorld(entityMap, cfg.StartingRoom)

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
