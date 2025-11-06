package dsl

import (
	"fmt"
	"strings"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/conditions"
)

type ConditionDef struct {
	Or *OrChain `parser:"@@"`
}

type OrChain struct {
	First *CondAtom     `parser:"@@"`
	Rest  []*OrChainRhs `parser:"( 'or' @@ )*"`
}

type OrChainRhs struct {
	Next *CondAtom `parser:"@@"`
}

type CondAtom struct {
	Paren      *ConditionDef             `parser:"  '(' @@ ')'"`
	Not        *NotCondition             `parser:"| @@"`
	Expr       *ExprCondition            `parser:"| @@"`
	HasTag     *HasTagCondition          `parser:"| @@"`
	IsPresent  *IsPresentCondition       `parser:"| @@"`
	RolesEqual *EventRolesEqualCondition `parser:"| @@"`
	HasChild   *HasChildCondition        `parser:"| @@"`
	MsgHas     *MessageContains          `parser:"| @@"`
}

type NotCondition struct {
	Cond *ConditionDef `parser:"'not' @@"`
}

type ExprCondition struct {
	Expr *Expression `parser:"'expr' '{' @@ '}'"`
}

type HasTagCondition struct {
	Target string `parser:"@Ident"`
	Tag    string `parser:"'has' 'tag' @String"`
}

type IsPresentCondition struct {
	Role string `parser:"@Ident 'exists'"`
}

type EventRolesEqualCondition struct {
	Role1 string `parser:"@Ident"`
	Role2 string `parser:"'is' @Ident"`
}

type HasChildCondition struct {
	ChildRole  string `parser:"@Ident"`
	ParentRole string `parser:"'in' @Ident"`
	Component  string `parser:"'.' @Ident"`
}

type MessageContains struct {
	Message string `parser:"'message' 'contains' @String"`
}

func (def *ConditionDef) Build() (entities.Condition, error) {
	if def == nil || def.Or == nil {
		return nil, fmt.Errorf("condition in when is empty")
	}

	acc, err := def.Or.First.Build()
	if err != nil {
		return nil, err
	}

	for _, rhs := range def.Or.Rest {
		next, err := rhs.Next.Build()
		if err != nil {
			return nil, err
		}
		acc = &conditions.Or{
			Left:  acc,
			Right: next,
		}
	}

	return acc, nil
}

func (def *CondAtom) Build() (entities.Condition, error) {
	if def == nil {
		return nil, fmt.Errorf("empty condition atom")
	}

	if def.Paren != nil {
		return def.Paren.Build()
	}

	if def.Not != nil {
		inner, err := def.Not.Cond.Build()
		if err != nil {
			return nil, fmt.Errorf("not condition: %w", err)
		}
		return &conditions.Not{Cond: inner}, nil
	}

	if def.Expr != nil {
		expression, err := def.Expr.Expr.Build()
		if err != nil {
			return nil, fmt.Errorf("condition expression: %w", err)
		}
		return &conditions.ExpressionTrue{Expression: expression}, nil
	}

	if def.HasTag != nil {
		eventRole, err := entities.ParseEventRole(def.HasTag.Target)
		if err != nil {
			return nil, fmt.Errorf("could not build has tag condition: %w", err)
		}
		return &conditions.HasTag{
			EventRole: eventRole,
			Tag:       def.HasTag.Tag,
		}, nil
	}

	if def.IsPresent != nil {
		eventRole, err := entities.ParseEventRole(def.IsPresent.Role)
		if err != nil {
			return nil, fmt.Errorf("could not build is-present condition: %w", err)
		}
		return &conditions.IsPresent{EventRole: eventRole}, nil
	}

	if def.RolesEqual != nil {
		role1, err := entities.ParseEventRole(def.RolesEqual.Role1)
		if err != nil {
			return nil, fmt.Errorf("event roles equal condition: %w", err)
		}
		role2, err := entities.ParseEventRole(def.RolesEqual.Role2)
		if err != nil {
			return nil, fmt.Errorf("event roles equal condition: %w", err)
		}
		return &conditions.EventRolesEqual{
			EventRole1: role1,
			EventRole2: role2,
		}, nil
	}

	if def.HasChild != nil {
		parentRole, err := entities.ParseEventRole(def.HasChild.ParentRole)
		if err != nil {
			return nil, fmt.Errorf("has child condition: %w", err)
		}
		component, err := entities.ParseComponentType(def.HasChild.Component)
		if err != nil {
			return nil, fmt.Errorf("has child condition: %w", err)
		}
		childRole, err := entities.ParseEventRole(def.HasChild.ChildRole)
		if err != nil {
			return nil, fmt.Errorf("has child condition: %w", err)
		}
		return &conditions.HasChild{
			ParentRole:    parentRole,
			ComponentType: component,
			ChildRole:     childRole,
		}, nil
	}

	if def.MsgHas != nil {
		return &conditions.MessageContains{
			MessageRegex: strings.ToLower(def.MsgHas.Message),
		}, nil
	}

	return nil, fmt.Errorf("unrecognized condition def")
}
