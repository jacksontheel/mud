package world

import (
	"fmt"

	"example.com/mud/world/entities"
)

type Player struct {
	world       *World
	entity      *entities.Entity
	currentRoom *entities.Room
}

func NewPlayer(name string, world *World, currentRoom *entities.Room) *Player {
	playerEntity := entities.NewEntity()
	playerEntity.Add(&entities.CDescriptioned{
		Text: fmt.Sprintf("%s the brave hero is here.", name),
	})
	playerEntity.Add(&entities.CAliased{
		Text: []string{name, "hero"},
	})
	playerEntity.Add(&entities.CTagged{
		Text: []string{"player"},
	})
	playerEntity.Add(&entities.CEventful{
		Rules: []entities.Rule{
			{
				When: entities.When{
					Type: "attack",
					Source: &entities.EntitySelector{
						Type: "self",
					},
				},
				Then: []entities.Action{
					&entities.ASay{
						Text: "Why are you hitting yourself?",
					},
				},
			},
		},
	})

	currentRoom.AddEntity(playerEntity)

	// TODO ERROR HANDLING
	return &Player{
		world:       world,
		entity:      playerEntity,
		currentRoom: currentRoom,
	}
}

func (p *Player) OpeningMessage() string {
	return fmt.Sprintf("You are a hero.\n%s", p.currentRoom.GetDescription(p.entity))
}

func (p *Player) Move(direction string) string {
	newRoom := p.world.GetNeighboringRoom(p.currentRoom, direction)
	if newRoom != nil {
		p.currentRoom.RemoveEntity(p.entity)
		p.currentRoom = newRoom
		p.currentRoom.AddEntity(p.entity)
		return p.currentRoom.GetDescription(p.entity)
	}
	return "You can't go there."
}

func (p *Player) Look(alias string) string {
	if alias == "" {
		return p.currentRoom.GetDescription(p.entity)
	}

	if target, ok := p.currentRoom.GetEntityByAlias(alias); ok {
		if eventful, ok := entities.Find[entities.Descriptioned](target); ok {
			return eventful.Description()
		}
		return fmt.Sprintf("The %s before you is undescribable.", alias)
	}
	return "you must be going mad. That's not here."

}

func (p *Player) Attack(alias string) string {
	return p.actUpon(
		"attack",
		alias,
		fmt.Sprintf("Now is not the time to attack %s.", alias),
	)
}

func (p *Player) Kiss(alias string) string {
	return p.actUpon(
		"kiss",
		alias,
		fmt.Sprintf("You can be romantic with %s later.", alias),
	)
}

func (p *Player) actUpon(action, alias, noMatchResponse string) string {
	if target, ok := p.currentRoom.GetEntityByAlias(alias); ok {
		if eventful, ok := entities.Find[entities.Eventful](target); ok {
			response, ok := eventful.OnEvent(&entities.Event{
				Type:   action,
				Source: p.entity,
				Target: target,
			})

			if ok {
				return response
			}
		}
		return noMatchResponse
	}
	return "you must be going mad. That's not here."
}
