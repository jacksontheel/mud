package dsl

import (
	"fmt"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/actions"
	"example.com/mud/world/entities/components"
)

type builtEntity struct {
	name string
	ent  *entities.Entity
	def  *EntityDef
}

// sidecar: prototypeName -> compName -> []childPrototypeNames
type childRefs map[string]map[string][]string

// turn a parsed DSL into concrete entities.
// returns a name->entity map.
// children references are materialized by recursively instantiating prototypes.
func buildAll(ast *DSL) (map[string]*entities.Entity, error) {
	if ast == nil {
		return nil, fmt.Errorf("nil DSL")
	}

	defs := make(map[string]*EntityDef, len(ast.Entities))
	for _, ed := range ast.Entities {
		if ed == nil {
			continue
		}
		if _, exists := defs[ed.Name]; exists {
			return nil, fmt.Errorf("duplicate entity %q", ed.Name)
		}
		defs[ed.Name] = ed
	}

	built := make(map[string]*builtEntity, len(defs))
	pending := make(childRefs, len(defs))

	for name, def := range defs {
		be, err := buildPrototype(name, def, pending)
		if err != nil {
			return nil, fmt.Errorf("build %s: %w", name, err)
		}
		built[name] = be
	}

	out := make(map[string]*entities.Entity, len(built))
	for name := range built {
		inst, err := instantiate(name, built, pending, map[string]bool{})
		if err != nil {
			return nil, fmt.Errorf("instantiate %s: %w", name, err)
		}
		out[name] = inst
	}
	return out, nil
}

// create component with components. collect child prototype names into the sidecar for later.
func buildPrototype(name string, def *EntityDef, pending childRefs) (*builtEntity, error) {
	e := entities.NewEntity()

	for _, block := range def.Blocks {
		if block.Rule != nil {
			eventful, ok := entities.GetComponent[*components.Eventful](e)
			if !ok {
				eventful = &components.Eventful{Rules: []*entities.Rule{}}
				e.Add(eventful)
			}
			rule, err := processRule(block.Rule)
			if err != nil {
				return nil, fmt.Errorf("could not process rule %s: %w", block.Rule.Command, err)
			}
			eventful.AddRule(rule)
			continue
		}

		if block.Component == nil {
			continue
		}

		comp, err := processComponentNoChildren(block.Component)
		if err != nil {
			return nil, fmt.Errorf("could not process component %s: %w", block.Component.Name, err)
		}
		e.Add(comp)

		for _, f := range block.Component.Fields {
			if f.Key == "children" {
				if pending[name] == nil {
					pending[name] = make(map[string][]string)
				}
				pending[name][block.Component.Name] =
					append(pending[name][block.Component.Name], f.Value.asStrings()...)
			}
		}
	}

	return &builtEntity{name: name, ent: e, def: def}, nil
}

// recursively instantiate a named prototype and wire up children for all child-holding components.
func instantiate(name string, protos map[string]*builtEntity, pending childRefs, visiting map[string]bool) (*entities.Entity, error) {
	be, ok := protos[name]
	if !ok {
		return nil, fmt.Errorf("unknown prototype %q", name)
	}
	if visiting[name] {
		return nil, fmt.Errorf("cycle detected at %q", name)
	}
	visiting[name] = true
	defer func() { visiting[name] = false }()

	inst := be.ent.Copy()

	// for each child-capable component on the entity, look up its pending child names from the prototype’s sidecar and attach recursively.
	if rm, ok := entities.GetComponent[*components.Room](inst); ok {
		slot := pending[name]["Room"]
		if len(slot) > 0 {
			for _, childName := range slot {
				childInst, err := instantiate(childName, protos, pending, visiting)
				if err != nil {
					return nil, err
				}
				rm.GetChildren().AddChild(childInst)
			}
		}
	}

	return inst, nil
}

func processRule(def *RuleDef) (*entities.Rule, error) {
	when, err := processWhen(def)
	if err != nil {
		return nil, fmt.Errorf("could not process 'when' for reaction on %s", def.Command)
	}

	then, err := processThen(def)
	if err != nil {
		return nil, fmt.Errorf("could not process 'then' for reaction on %s", def.Command)
	}

	return &entities.Rule{
		When: when,
		Then: then,
	}, nil
}

func processWhen(def *RuleDef) (*entities.When, error) {
	return &entities.When{
		Type: def.Command,
	}, nil
}

func processThen(def *RuleDef) ([]entities.Action, error) {
	ret := make([]entities.Action, len(def.Actions))

	for i, aDef := range def.Actions {
		var newAction entities.Action

		if aDef.Say != nil {
			newAction = &actions.Say{
				Text: aDef.Say.Value,
			}
		}

		ret[i] = newAction
	}

	return ret, nil
}

// turn DSL pairs into map[string]string map
func (l *Literal) asMap() map[string]string {
	if len(l.Pairs) == 0 {
		return nil
	}
	m := make(map[string]string, len(l.Pairs))
	for _, kv := range l.Pairs {
		m[kv.Key] = kv.Value
	}
	return m
}

// turn DSL Strings into []string
func (l *Literal) asStrings() []string {
	return append([]string(nil), l.Strings...)
}
