package dsl

import (
	"fmt"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/actions"
	"example.com/mud/world/entities/components"
)

type collectedDefs struct {
	entitiesById map[string]*EntityDef
	traitsById   map[string]*TraitDef
}

type ChildrenPlan map[string]map[entities.ComponentType][]string

type LoweredEntity struct {
	name        string
	description string
	tags        []string
	aliases     []string
	components  []entities.Component
	rules       []*entities.Rule
}

type entityPrototype struct {
	id  string
	ent *entities.Entity
	def *EntityDef
}

type entityPrototypes struct {
	prototypesById map[string]*entityPrototype
	traitsById     map[string]*TraitDef
	childrenPlan   ChildrenPlan
	visiting       map[string]struct{}
}

func Compile(ast *DSL) (map[string]*entities.Entity, error) {
	if ast == nil {
		return nil, fmt.Errorf("nil DSL")
	}

	collectedDefs, err := collectDefs(ast.Declarations)
	if err != nil {
		return nil, fmt.Errorf("could not collect top level declarations: %w", err)
	}

	prototypes, err := collectedDefs.collectPrototypes()
	if err != nil {
		return nil, fmt.Errorf("could not collect prototype entities: %w", err)
	}

	entitiesById, err := prototypes.instantiatePrototypes()
	if err != nil {
		return nil, fmt.Errorf("could not instantiate prototype entities: %w", err)
	}

	return entitiesById, nil
}

// collect entity and trait definitions
func collectDefs(decls []*TopLevel) (*collectedDefs, error) {
	entitiesById := make(map[string]*EntityDef, len(decls))
	traitsById := make(map[string]*TraitDef, len(decls))

	for _, declaration := range decls {
		if declaration == nil {
			return nil, fmt.Errorf("declaration at top level is nil")
		}

		if ed := declaration.Entity; ed != nil {
			if _, exists := entitiesById[ed.Name]; exists {
				return nil, fmt.Errorf("duplicate entity %s", ed.Name)
			}

			entitiesById[ed.Name] = ed
		} else if td := declaration.Trait; td != nil {
			if _, exists := entitiesById[td.Name]; exists {
				return nil, fmt.Errorf("duplicate trait %s", td.Name)
			}

			traitsById[declaration.Trait.Name] = declaration.Trait
		} else {
			return nil, fmt.Errorf("declaration at top level is empty")
		}
	}

	return &collectedDefs{
		entitiesById: entitiesById,
		traitsById:   traitsById,
	}, nil
}

// expand traits in each entity definition
func (c *collectedDefs) collectPrototypes() (*entityPrototypes, error) {
	ep := &entityPrototypes{
		prototypesById: map[string]*entityPrototype{},
		traitsById:     c.traitsById,
		childrenPlan:   map[string]map[entities.ComponentType][]string{},
		visiting:       map[string]struct{}{},
	}

	// build prototypes of each entity and put them in name->builtEntity map
	for name, ed := range c.entitiesById {
		// build prototype and populate pending children
		prototypeEntity, err := ep.buildPrototype(name, ed.Blocks)
		if err != nil {
			return nil, fmt.Errorf("build %s: %w", name, err)
		}
		ep.prototypesById[name] = &entityPrototype{
			id:  name,
			ent: prototypeEntity,
			def: ed,
		}
	}

	return ep, nil
}

// create prototype entity with components. collect child prototype names into the sidecar for later.
func (ep *entityPrototypes) buildPrototype(id string, blocks []*EntityBlock) (*entities.Entity, error) {

	loweredEntity, err := ep.lowerEntity(id, blocks)
	if err != nil {
		return nil, fmt.Errorf("could not build prototype: %w", err)
	}

	e := entities.NewEntity(
		loweredEntity.name,
		loweredEntity.description,
		loweredEntity.aliases,
		loweredEntity.tags,
		nil,
	)

	for _, c := range loweredEntity.components {
		e.Add(c)
	}

	if len(loweredEntity.rules) > 0 {
		// create eventful if it isn't already
		eventful, ok := entities.GetComponent[*components.Eventful](e)
		if !ok {
			eventful = &components.Eventful{Rules: []*entities.Rule{}}
			e.Add(eventful)
		}

		for _, r := range loweredEntity.rules {
			eventful.AddRule(r)
		}
	}

	for _, block := range blocks {
		if block.Component == nil {
			continue
		}

		for _, f := range block.Component.Fields {
			if f.Key == "children" {
				if ep.childrenPlan[id] == nil {
					ep.childrenPlan[id] = make(map[entities.ComponentType][]string)
				}

				// populate pending children map
				componentType, err := entities.ParseComponentType(block.Component.Name)
				if err != nil {
					return nil, fmt.Errorf("could not build prototype '%s': %w", id, err)
				}

				ep.childrenPlan[id][componentType] =
					append(ep.childrenPlan[id][componentType], f.Value.Strings...)
			}
		}
	}

	return e, nil
}

// recursively expand traits in entities
func (ep *entityPrototypes) lowerEntity(id string, blocks []*EntityBlock) (*LoweredEntity, error) {
	if _, ok := ep.visiting[id]; ok {
		return nil, fmt.Errorf("cycle detected at %q", id)
	}
	ep.visiting[id] = struct{}{}
	defer func() { delete(ep.visiting, id) }()

	var name string
	var description string
	var aliases []string
	var tags []string

	components := make([]entities.Component, 0, len(blocks))
	rules := make([]*entities.Rule, 0, len(blocks))

	for _, block := range blocks {
		// process rules
		if block.Rule != nil {
			rule, err := buildRule(block.Rule)
			if err != nil {
				return nil, fmt.Errorf("could not process rule %s: %w", block.Rule.Command, err)
			}
			rules = append(rules, rule)
		} else if block.Component != nil {
			// process component into prototype without children
			comp, err := processComponentPrototype(block.Component)
			if err != nil {
				return nil, fmt.Errorf("could not process component %s: %w", block.Component.Name, err)
			}
			components = append(components, comp)
		} else if block.Trait != nil {
			loweredTrait, err := ep.lowerEntity(block.Trait.Name, ep.traitsById[block.Trait.Name].Blocks)
			if err != nil {
				return nil, fmt.Errorf("could not process trait '%s': %w", block.Trait.Name, err)
			}

			components = append(components, loweredTrait.components...)
			rules = append(rules, loweredTrait.rules...)
		} else if block.Field != nil {
			f := block.Field
			switch f.Key {
			case "name":
				if f.Value.String == nil {
					return nil, fmt.Errorf("identity.name must be a string")
				}
				name = *f.Value.String
			case "description":
				if f.Value.String == nil {
					return nil, fmt.Errorf("identity.description must be a string")
				}
				description = *f.Value.String
			case "aliases":
				aliases = f.Value.Strings
			case "tags":
				tags = f.Value.Strings
			default:
				return nil, fmt.Errorf("identity: unknown field %s", f.Key)
			}
		} else {
			return nil, fmt.Errorf("could not expand empty entity block")
		}
	}

	// only do verification if at a top level entity
	if len(ep.visiting) == 1 {
		// verify name, description, and aliases are set. Empty tags is ok
		if name == "" {
			return nil, fmt.Errorf("entity '%s' has no name", id)
		}
		if description == "" {
			return nil, fmt.Errorf("entity '%s' has no description", id)
		}
		if len(aliases) == 0 {
			return nil, fmt.Errorf("entity '%s' has no aliases", id)
		}
	}

	return &LoweredEntity{
		name:        name,
		description: description,
		tags:        tags,
		aliases:     aliases,
		components:  components,
		rules:       rules,
	}, nil
}

// loop through prototypes and instantiate them into a map of entities by name
func (ep *entityPrototypes) instantiatePrototypes() (map[string]*entities.Entity, error) {
	out := make(map[string]*entities.Entity, len(ep.prototypesById))
	for name := range ep.prototypesById {
		entity, err := ep.instantiate(name, nil)
		if err != nil {
			return nil, fmt.Errorf("could not instantiate '%s': %w", name, err)
		}
		out[name] = entity
	}
	return out, nil
}

// recursively instantiate a named prototype and wire up children for all child-holding components.
func (ep *entityPrototypes) instantiate(id string, parent entities.ComponentWithChildren) (*entities.Entity, error) {
	be, ok := ep.prototypesById[id]
	if !ok {
		return nil, fmt.Errorf("unknown prototype %q", id)
	}
	if _, ok := ep.visiting[id]; ok {
		return nil, fmt.Errorf("cycle detected at %q", id)
	}
	ep.visiting[id] = struct{}{}
	defer func() { delete(ep.visiting, id) }()

	inst := be.ent.Copy(parent)

	// for each child-capable component on the entity, look up its pending child names from the prototype’s sidecar and attach recursively.
	if rm, ok := entities.GetComponent[*components.Room](inst); ok {
		slot := ep.childrenPlan[id][entities.ComponentRoom]
		if len(slot) > 0 {
			for _, childName := range slot {
				childInst, err := ep.instantiate(childName, rm)
				if err != nil {
					return nil, err
				}
				rm.GetChildren().AddChild(childInst)
			}
		}
	}

	return inst, nil
}

func buildRule(def *RuleDef) (*entities.Rule, error) {
	when, err := buildWhen(def)
	if err != nil {
		return nil, fmt.Errorf("could not process 'when' for reaction on %s", def.Command)
	}

	then, err := buildThen(def)
	if err != nil {
		return nil, fmt.Errorf("could not process 'then' for reaction on %s: %w", def.Command, err)
	}

	return &entities.Rule{
		When: when,
		Then: then,
	}, nil
}

func buildWhen(def *RuleDef) (*entities.When, error) {
	sourceSelector, err := buildEntitySelector(def.By)
	if err != nil {
		return nil, fmt.Errorf("could not create selector for when reference %s: %w", *def.By, err)
	}

	instrumentSelector, err := buildEntitySelector(def.With)
	if err != nil {
		return nil, fmt.Errorf("could not create selector for when reference %s: %w", *def.With, err)
	}

	return &entities.When{
		Type:       def.Command,
		Source:     sourceSelector,
		Instrument: instrumentSelector,
	}, nil
}

func buildEntitySelector(ref *string) (*entities.EntitySelector, error) {
	if ref == nil || *ref == "" {
		return nil, nil
	}
	value := *ref

	switch value[0] {
	case '#': // tag
		value = value[1:]
		return &entities.EntitySelector{
			Type:  "tag",
			Value: value,
		}, nil
	default:
		return nil, fmt.Errorf("illegal value reference for when: %s", value)
	}
}

func buildThen(def *RuleDef) ([]entities.Action, error) {
	ret := make([]entities.Action, len(def.Actions))

	for i, aDef := range def.Actions {
		var newAction entities.Action

		if aDef.Print != nil {
			printTarget := actions.StringToEventRole(aDef.Print.Target)
			if printTarget == actions.EventRoleUnknown {
				return nil, fmt.Errorf("unknown print target %s", aDef.Print.Target)
			}

			newAction = &actions.Print{
				Text:      aDef.Print.Value,
				EventRole: actions.StringToEventRole(aDef.Print.Target),
			}
		} else if aDef.Publish != nil {
			newAction = &actions.Publish{
				Text: aDef.Publish.Value,
			}
		} else if aDef.Copy != nil {
			copyTarget := actions.StringToEventRole(aDef.Copy.Target)
			if copyTarget == actions.EventRoleUnknown {
				return nil, fmt.Errorf("unknown copy target %s", aDef.Copy.Target)
			}

			component, err := entities.ParseComponentType(aDef.Copy.Component)
			if err != nil {
				return nil, fmt.Errorf("could not build action: %w", err)
			}

			newAction = &actions.Copy{
				EntityId:      aDef.Copy.EntityId,
				EventRole:     actions.StringToEventRole(aDef.Copy.Target),
				ComponentType: component,
			}
		} else if aDef.Move != nil {
			roleOrigin := actions.StringToEventRole(aDef.Move.RoleOrigin)
			if roleOrigin == actions.EventRoleUnknown {
				return nil, fmt.Errorf("unknown move origin role %s", aDef.Copy.Target)
			}

			roleDestination := actions.StringToEventRole(aDef.Move.RoleDestination)
			if roleDestination == actions.EventRoleUnknown {
				return nil, fmt.Errorf("unknown move destination role %s", aDef.Copy.Target)
			}

			component, err := entities.ParseComponentType(aDef.Move.Component)
			if err != nil {
				return nil, fmt.Errorf("could not build action: %w", err)
			}

			newAction = &actions.Move{
				RoleOrigin:      roleOrigin,
				RoleDestination: roleDestination,
				ComponentType:   component,
			}
		} else {
			return nil, fmt.Errorf("action in then is empty")
		}

		ret[i] = newAction
	}

	return ret, nil
}
