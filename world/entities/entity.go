package entities

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"example.com/mud/models"
	"example.com/mud/utils"
)

type Entity struct {
	mu         sync.RWMutex
	components map[reflect.Type]Component

	Name        string
	Description string
	Aliases     []string
	Tags        []string
	Fields      map[string]any
	Parent      ComponentWithChildren
}

func NewEntity(name, description string, aliases []string, tags []string, fields map[string]any, parent ComponentWithChildren) *Entity {
	return &Entity{
		components:  map[reflect.Type]Component{},
		Name:        name,
		Description: description,
		Aliases:     aliases,
		Tags:        tags,
		Fields:      fields,
		Parent:      parent,
	}
}

func (e *Entity) Copy(parent ComponentWithChildren) *Entity {
	newEntity := NewEntity(e.Name, e.Description, e.Aliases, e.Tags, e.Fields, parent)
	for _, c := range e.components {
		newEntity.Add(c.Copy())
	}
	return newEntity
}

func (e *Entity) Add(c Component) *Entity {
	e.mu.Lock()
	e.components[reflect.TypeOf(c)] = c
	e.mu.Unlock()
	return e
}

func (e *Entity) GetComponentWithChildren(ct ComponentType) (ComponentWithChildren, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, c := range e.components {
		id := c.Id()
		if cwc, ok := any(c).(ComponentWithChildren); ok {
			if id == ct {
				return cwc, true
			}
		}
	}

	return nil, false
}

func (e *Entity) RequireComponentWithChildren(ct ComponentType) (ComponentWithChildren, error) {
	c, ok := e.GetComponentWithChildren(ct)

	if !ok {
		return nil, fmt.Errorf("entity does not have component with children %s", ct.String())
	}

	return c, nil
}

func (e *Entity) GetComponentsWithChildren() []ComponentWithChildren {
	e.mu.RLock()
	defer e.mu.RUnlock()

	components := make([]ComponentWithChildren, 0, len(e.components))

	for _, c := range e.components {
		if cwc, ok := any(c).(ComponentWithChildren); ok {
			components = append(components, cwc)
		}
	}

	return components
}

func (e *Entity) GetDescription() (string, error) {
	var b strings.Builder

	formatted, err := utils.FormatText(e.Description, map[string]string{})
	if err != nil {
		return "", fmt.Errorf("could not format description for entity '%s': %w", e.Name, err)
	}

	b.WriteString(fmt.Sprintf("- %s", formatted))

	for _, cwc := range e.GetComponentsWithChildren() {
		if !cwc.GetChildren().GetRevealed() {
			continue
		}

		children := cwc.GetChildren().GetChildren()
		if len(children) == 0 {
			continue
		}

		var childB strings.Builder
		childB.WriteString("\n")

		childB.WriteString(fmt.Sprintf("%s%s", models.Tab, cwc.GetChildren().GetPrefix()))
		childB.WriteString(" (\n")

		for _, child := range children {
			cDescription, err := child.GetDescription()
			if err != nil {
				return "", fmt.Errorf("could not format description for entity '%s': %w", child.Name, err)
			}

			childB.WriteString(fmt.Sprintf("%s%s%s", models.Tab, models.Tab, cDescription))
			childB.WriteString("\n")
		}

		b.WriteString(childB.String())

		b.WriteString(fmt.Sprintf("%s)", models.Tab))
	}

	return b.String(), nil
}

func GetComponent[C Component](e *Entity) (C, bool) {
	var zero C
	e.mu.RLock()
	defer e.mu.RUnlock()
	v, ok := e.components[reflect.TypeOf((*C)(nil)).Elem()]
	if !ok {
		return zero, false
	}
	return v.(C), true
}

func RequireComponent[C Component](e *Entity) (C, error) {
	c, ok := GetComponent[C](e)

	if ok {
		return c, nil
	}

	var zero C
	return zero, fmt.Errorf("entity does not have component %s", zero.Id().String())
}
