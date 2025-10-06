package world

import (
	"fmt"
	"strings"

	"example.com/mud/models"
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

func (p *Player) Map() (string, error) {
	coordByRoom, err := assignCoordinates(p.currentRoom, p.world, 5)
	if err != nil {
		return "", fmt.Errorf("map: assign coordinates: %w", err)
	}

	currentRoom, err := entities.RequireComponent[*components.Room](p.currentRoom)
	if err != nil {
		return "", fmt.Errorf("cannot map non-room area: %w", err)
	}

	ascii, err := renderMap(coordByRoom, currentRoom, p.world)
	if err != nil {
		return "", fmt.Errorf("map: render ascii: %w", err)
	}

	return ascii, nil
}

type coord struct{ X, Y int }

var cardinalDelta = map[string]coord{
	"north": {0, -1},
	"south": {0, 1},
	"east":  {1, 0},
	"west":  {-1, 0},
}

type mapper struct {
	coordByRoom map[*components.Room]coord
	roomAtCoord map[coord]*components.Room
	maxDepth    int
	world       *World
}

func assignCoordinates(start *entities.Entity, world *World, maxDepth int) (map[*components.Room]coord, error) {
	m := &mapper{
		coordByRoom: make(map[*components.Room]coord),
		roomAtCoord: make(map[coord]*components.Room),
		maxDepth:    maxDepth,
		world:       world,
	}

	err := m.visit(start, 0, 0, 0)
	return m.coordByRoom, err
}

func (m *mapper) visit(e *entities.Entity, x, y, depth int) error {
	if depth > m.maxDepth {
		return nil
	}

	// require visited entity to be a room
	r, err := entities.RequireComponent[*components.Room](e)
	if err != nil {
		return fmt.Errorf("cannot map non-room: %w", err)
	}

	coord := coord{x, y}
	if existing, exists := m.roomAtCoord[coord]; exists && existing != r {
		return fmt.Errorf("mapping non-euclidian space: %w", err)
	}
	if _, visited := m.coordByRoom[r]; visited {
		return nil
	}

	m.coordByRoom[r] = coord
	m.roomAtCoord[coord] = r

	for direction, roomId := range r.GetExits() {
		delta, ok := cardinalDelta[direction]
		if !ok {
			// we don't map up or down
			continue
		}

		nextRoomEntity, ok := m.world.entityMap[roomId]
		if !ok {
			return fmt.Errorf("entity with id '%s' does not exist: %w", roomId, err)
		}

		if err := m.visit(nextRoomEntity, x+delta.X, y+delta.Y, depth+1); err != nil {
			return fmt.Errorf("recursive visit: %w", err)
		}
	}
	return nil
}

func renderMap(coordByRoom map[*components.Room]coord, currentRoom *components.Room, world *World) (string, error) {
	if len(coordByRoom) == 0 {
		return "", nil
	}

	minX, maxX, minY, maxY := 0, 0, 0, 0
	for _, p := range coordByRoom {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	width := (maxX-minX)*2 + 1
	height := (maxY-minY)*2 + 1
	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	for r, c := range coordByRoom {
		gx := (c.X - minX) * 2
		gy := (c.Y - minY) * 2

		if r == currentRoom {
			grid[gy][gx] = fmt.Sprintf("%s%s%s", models.SGR["red"], "@", models.SGR["reset"])
		} else {
			grid[gy][gx] = "O"
		}

		for _, roomId := range r.GetExits() {
			roomEntity, ok := world.entityMap[roomId]
			if !ok {
				return "", fmt.Errorf("entity with id '%s' does not exist", roomId)
			}

			room, err := entities.RequireComponent[*components.Room](roomEntity)
			if err != nil {
				return "", fmt.Errorf("render map: %w", err)
			}

			npos, ok := coordByRoom[room]
			if !ok {
				continue
			}
			dx := npos.X - c.X
			dy := npos.Y - c.Y
			if dx == 1 {
				grid[gy][gx+1] = "-"
			} else if dx == -1 {
				grid[gy][gx-1] = "-"
			} else if dy == 1 {
				grid[gy+1][gx] = "|"
			} else if dy == -1 {
				grid[gy-1][gx] = "|"
			}
		}
	}

	var b strings.Builder
	b.Grow(height * (width + 1))
	for _, row := range grid {
		b.WriteString(strings.Join(row, ""))
		b.WriteByte('\n')
	}
	return b.String(), nil
}
