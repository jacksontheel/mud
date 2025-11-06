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

func (ad *ActionDef) BuildAction() (entities.Action, error) {
	switch {
	case ad.Print != nil:
		eventRole, err := entities.ParseEventRole(ad.Print.Target)
		if err != nil {
			return nil, fmt.Errorf("could not build print action: %w", err)
		}

		return &actions.Print{
			Text:      ad.Print.Value,
			EventRole: eventRole,
		}, nil
	case ad.Publish != nil:
		return &actions.Publish{
			Text: ad.Publish.Value,
		}, nil
	case ad.Copy != nil:
		eventRole, err := entities.ParseEventRole(ad.Copy.Target)
		if eventRole == entities.EventRoleUnknown {
			return nil, fmt.Errorf("could not build copy action: %w", err)
		}

		component, err := entities.ParseComponentType(ad.Copy.Component)
		if err != nil {
			return nil, fmt.Errorf("could not build action: %w", err)
		}

		return &actions.Copy{
			EntityId:      ad.Copy.EntityId,
			EventRole:     eventRole,
			ComponentType: component,
		}, nil
	case ad.Move != nil:
		roleObject, err := entities.ParseEventRole(ad.Move.RoleObject)
		if err != nil {
			return nil, fmt.Errorf("could not build move action for origin: %w", err)
		}

		roleDestination, err := entities.ParseEventRole(ad.Move.RoleDestination)
		if err != nil {
			return nil, fmt.Errorf("could not build move action for destination: %w", err)
		}

		component, err := entities.ParseComponentType(ad.Move.Component)
		if err != nil {
			return nil, fmt.Errorf("could not build action: %w", err)
		}

		return &actions.Move{
			RoleObject:      roleObject,
			RoleDestination: roleDestination,
			ComponentType:   component,
		}, nil
	case ad.SetField != nil:
		role, err := entities.ParseEventRole(ad.SetField.Role)
		if err != nil {
			return nil, fmt.Errorf("event set field action: %w", err)
		}

		expression, err := ad.SetField.Expr.Build()
		if err != nil {
			return nil, fmt.Errorf("expression set field action: %w", err)
		}

		return &actions.SetField{
			Role:       role,
			Field:      ad.SetField.Field,
			Expression: expression,
		}, nil
	case ad.DestroyAction != nil:
		role, err := entities.ParseEventRole(ad.DestroyAction.Role)
		if err != nil {
			return nil, fmt.Errorf("event destroy action: %w", err)
		}

		return &actions.Destroy{
			Role: role,
		}, nil
	case ad.RevealChildrenAction != nil:
		role, err := entities.ParseEventRole(ad.RevealChildrenAction.Role)
		if err != nil {
			return nil, fmt.Errorf("could not build reveal children action for role: %w", err)
		}

		component, err := entities.ParseComponentType(ad.RevealChildrenAction.Component)
		if err != nil {
			return nil, fmt.Errorf("could not build reveal children action: %w", err)
		}

		return &actions.RevealChildren{
			Role:          role,
			ComponentType: component,
			Reveal:        ad.RevealChildrenAction.Set == "reveal",
		}, nil
	}

	return nil, fmt.Errorf("action is empty")
}
