// ir/entity.go
package ir

import "example.com/mud/world/entities"

type Entity struct {
	ID          string
	Name        string
	Description string
	Aliases     []string
	Tags        []string
	Components  Components
}

func (e *Entity) Build() *entities.Entity {
	ent := entities.NewEntity(
		e.Name,
		e.Description,
		e.Aliases,
		e.Tags,
		nil,
		nil,
	)

	components := e.Components.Build()
	for _, c := range components {
		ent.Add(c)
	}

	// If this entity is a room and has inline children, build and attach them now.
	if e.Components.Room != nil && len(e.Components.Room.InlineChildren) > 0 {
		roomComp, err := ent.RequireComponentWithChildren(entities.ComponentRoom)
		if err == nil {
			children := roomComp.GetChildren()
			for _, childIR := range e.Components.Room.InlineChildren {
				childRT := childIR.Build() // recursive build
				children.AddChild(childRT) // attach directly; no second pass
			}
		}
	}

	return ent
}
