package dsl

import (
	"fmt"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/actions"
)

type ActionDef struct {
	Print                *PrintAction          `parser:"  'print' @@"`
	Publish              *PublishAction        `parser:"| 'publish' @@"`
	Copy                 *CopyAction           `parser:"| 'copy' @@"`
	Move                 *MoveAction           `parser:"| 'move' @@"`
	SetField             *SetFieldAction       `parser:"| 'set' @@"`
	DestroyAction        *DestroyAction        `parser:"| 'destroy' @@"`
	RevealChildrenAction *RevealChildrenAction `parser:"| @@"`
}

type PrintAction struct {
	Target string `parser:"@Ident"`
	Value  string `parser:"@String"`
}

type PublishAction struct {
	Value string `parser:"@String"`
}

type CopyAction struct {
	EntityId  string `parser:"@String"`
	Target    string `parser:"'to' @Ident"`
	Component string `parser:"'.' @Ident"`
}

type MoveAction struct {
	RoleObject      string `parser:"@Ident"`
	RoleDestination string `parser:"'to' @Ident"`
	Component       string `parser:"'.' @Ident"`
}

type SetFieldAction struct {
	Role  string     `parser:"@Ident"`
	Field string     `parser:"'.' @Ident"`
	Expr  Expression `parser:"'to' @@"`
}

type RevealChildrenAction struct {
	Set       string `parser:"@('reveal' | 'hide')"`
	Role      string `parser:"@Ident"`
	Component string `parser:"'.' @Ident"`
}

type DestroyAction struct {
	Role string `parser:"@Ident"`
}

func (def *ActionDef) Build() (entities.Action, error) {
	switch {
	case def.Print != nil:
		return def.Print.Build()
	case def.Publish != nil:
		return def.Publish.Build()
	case def.Copy != nil:
		return def.Copy.Build()
	case def.Move != nil:
		return def.Move.Build()
	case def.SetField != nil:
		return def.SetField.Build()
	case def.DestroyAction != nil:
		return def.DestroyAction.Build()
	case def.RevealChildrenAction != nil:
		return def.RevealChildrenAction.Build()
	}

	return nil, fmt.Errorf("action is empty")
}

func (def *PrintAction) Build() (entities.Action, error) {
	eventRole, err := entities.ParseEventRole(def.Target)
	if err != nil {
		return nil, fmt.Errorf("could not build print action: %w", err)
	}

	return &actions.Print{
		Text:      def.Value,
		EventRole: eventRole,
	}, nil
}

func (def *PublishAction) Build() (entities.Action, error) {
	return &actions.Publish{
		Text: def.Value,
	}, nil
}

func (def *CopyAction) Build() (entities.Action, error) {
	eventRole, err := entities.ParseEventRole(def.Target)
	if eventRole == entities.EventRoleUnknown {
		return nil, fmt.Errorf("could not build copy action: %w", err)
	}

	component, err := entities.ParseComponentType(def.Component)
	if err != nil {
		return nil, fmt.Errorf("could not build action: %w", err)
	}

	return &actions.Copy{
		EntityId:      def.EntityId,
		EventRole:     eventRole,
		ComponentType: component,
	}, nil
}

func (def *MoveAction) Build() (entities.Action, error) {
	roleObject, err := entities.ParseEventRole(def.RoleObject)
	if err != nil {
		return nil, fmt.Errorf("could not build move action for origin: %w", err)
	}

	roleDestination, err := entities.ParseEventRole(def.RoleDestination)
	if err != nil {
		return nil, fmt.Errorf("could not build move action for destination: %w", err)
	}

	component, err := entities.ParseComponentType(def.Component)
	if err != nil {
		return nil, fmt.Errorf("could not build action: %w", err)
	}

	return &actions.Move{
		RoleObject:      roleObject,
		RoleDestination: roleDestination,
		ComponentType:   component,
	}, nil
}

func (def *SetFieldAction) Build() (entities.Action, error) {
	role, err := entities.ParseEventRole(def.Role)
	if err != nil {
		return nil, fmt.Errorf("event set field action: %w", err)
	}

	expression, err := def.Expr.Build()
	if err != nil {
		return nil, fmt.Errorf("expression set field action: %w", err)
	}

	return &actions.SetField{
		Role:       role,
		Field:      def.Field,
		Expression: expression,
	}, nil
}

func (def *DestroyAction) Build() (entities.Action, error) {
	role, err := entities.ParseEventRole(def.Role)
	if err != nil {
		return nil, fmt.Errorf("event destroy action: %w", err)
	}

	return &actions.Destroy{
		Role: role,
	}, nil
}

func (def *RevealChildrenAction) Build() (entities.Action, error) {
	role, err := entities.ParseEventRole(def.Role)
	if err != nil {
		return nil, fmt.Errorf("could not build reveal children action for role: %w", err)
	}

	component, err := entities.ParseComponentType(def.Component)
	if err != nil {
		return nil, fmt.Errorf("could not build reveal children action: %w", err)
	}

	return &actions.RevealChildren{
		Role:          role,
		ComponentType: component,
		Reveal:        def.Set == "reveal",
	}, nil
}
