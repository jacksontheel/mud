package dsl

import (
	"fmt"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
)

// Optional helpers on your Literal node.
func (l *Literal) AsMap() map[string]string {
	if len(l.Pairs) == 0 {
		return nil
	}
	m := make(map[string]string, len(l.Pairs))
	for _, kv := range l.Pairs {
		m[kv.Key] = kv.Value
	}
	return m
}

func (l *Literal) AsStrings() []string {
	return append([]string(nil), l.Strings...)
}

type builtEntity struct {
	ent             *entities.Entity
	pendingChildren []string // names to resolve after all prototypes exist
}

// BuildAll compiles a parsed DSL into concrete runtime entities.
// It returns a name->entity map. "children" references are materialized as CLONED entities.
func BuildAll(ast *DSL) (map[string]*entities.Entity, error) {
	if ast == nil {
		return nil, fmt.Errorf("nil DSL")
	}

	// 1) Index defs by name for error messages, etc.
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

	// 2) First pass: build prototypes (no children resolved yet).
	built := make(map[string]*builtEntity, len(defs))
	for name, def := range defs {
		be, err := buildPrototype(def)
		if err != nil {
			return nil, fmt.Errorf("build %s: %w", name, err)
		}
		built[name] = be
	}

	// 3) Second pass: resolve children by cloning referenced prototypes.
	for name, be := range built {
		if len(be.pendingChildren) == 0 {
			continue
		}
		room, ok := entities.GetComponent[*components.Room](be.ent)
		if !ok {
			return nil, fmt.Errorf("entity %q has children but no Room component", name)
		}
		for _, childName := range be.pendingChildren {
			proto, ok := built[childName]
			if !ok {
				return nil, fmt.Errorf("entity %q references unknown child %q", name, childName)
			}
			room.GetChildren().AddChild(proto.ent.Copy())
		}
	}

	// 4) Return final map
	out := make(map[string]*entities.Entity, len(built))
	for name, be := range built {
		out[name] = be.ent
	}
	return out, nil
}

// buildPrototype creates an entity and fills components. It collects "children" names to resolve later.
func buildPrototype(def *EntityDef) (*builtEntity, error) {
	e := entities.NewEntity()
	var pending []string

	for _, comp := range def.Components {
		switch comp.Name {
		case "Identity":
			id := &components.Identity{}
			for _, f := range comp.Fields {
				switch f.Key {
				case "name":
					if f.Value.String == nil {
						return nil, fmt.Errorf("Identity.name must be a string")
					}
					id.Name = *f.Value.String
				case "description":
					if f.Value.String == nil {
						return nil, fmt.Errorf("Identity.description must be a string")
					}
					id.Description = *f.Value.String
				case "aliases":
					id.Aliases = f.Value.AsStrings()
				case "tags":
					id.Tags = f.Value.AsStrings()
				default:
					return nil, fmt.Errorf("Identity: unknown field %q", f.Key)
				}
			}
			e.Add(id)

		case "Room":
			rm := components.NewRoom()
			for _, f := range comp.Fields {
				switch f.Key {
				case "exits":
					m := f.Value.AsMap()
					if m == nil {
						m = map[string]string{}
					}
					rm.SetExits(m)
				case "children":
					pending = append(pending, f.Value.AsStrings()...)
				default:
					return nil, fmt.Errorf("Room: unknown field %q", f.Key)
				}
			}
			e.Add(rm)

		default:
			return nil, fmt.Errorf("unknown component %q", comp.Name)
		}
	}

	return &builtEntity{ent: e, pendingChildren: pending}, nil
}
