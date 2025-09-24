package dsl

import (
	"fmt"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/actions"
	"example.com/mud/world/entities/components"
)

type collectedDefs struct {
	entitiesByName map[string]*EntityDef
	traitsByName   map[string]*TraitDef
}

type ChildrenPlan map[string]map[entities.ComponentType][]string

type entityPrototype struct {
	name string
	ent  *entities.Entity
	def  *EntityDef
}

type entityPrototypes struct {
	prototypesByName map[string]*entityPrototype
	traitsByName     map[string]*TraitDef
	childrenPlan     ChildrenPlan
	visiting         map[string]struct{}
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

	entitiesByName, err := prototypes.instantiatePrototypes()
	if err != nil {
		return nil, fmt.Errorf("could not instantiate prototype entities: %w", err)
	}

	return entitiesByName, nil
}

// collect entity and trait definitions
func collectDefs(decls []*TopLevel) (*collectedDefs, error) {
	entitiesByName := make(map[string]*EntityDef, len(decls))
	traitsByName := make(map[string]*TraitDef, len(decls))

	for _, declaration := range decls {
		if declaration == nil {
			return nil, fmt.Errorf("declaration at top level is nil")
		}

		if ed := declaration.Entity; ed != nil {
			if _, exists := entitiesByName[ed.Name]; exists {
				return nil, fmt.Errorf("duplicate entity %s", ed.Name)
			}

			entitiesByName[ed.Name] = ed
		} else if td := declaration.Trait; td != nil {
			if _, exists := entitiesByName[td.Name]; exists {
				return nil, fmt.Errorf("duplicate trait %s", td.Name)
			}

			traitsByName[declaration.Trait.Name] = declaration.Trait
		} else {
			return nil, fmt.Errorf("declaration at top level is empty")
		}
	}

	return &collectedDefs{
		entitiesByName: entitiesByName,
		traitsByName:   traitsByName,
	}, nil
}

// expand traits in each entity definition
func (c *collectedDefs) collectPrototypes() (*entityPrototypes, error) {
	ep := &entityPrototypes{
		prototypesByName: map[string]*entityPrototype{},
		traitsByName:     c.traitsByName,
		childrenPlan:     map[string]map[entities.ComponentType][]string{},
		visiting:         map[string]struct{}{},
	}

	// build prototypes of each entity and put them in name->builtEntity map
	for name, ed := range c.entitiesByName {
		// build prototype and populate pending children
		prototypeEntity, err := ep.buildPrototype(name, ed.Blocks)
		if err != nil {
			return nil, fmt.Errorf("build %s: %w", name, err)
		}
		ep.prototypesByName[name] = &entityPrototype{
			name: name,
			ent:  prototypeEntity,
			def:  ed,
		}
	}

	return ep, nil
}

// create prototype entity with components. collect child prototype names into the sidecar for later.
func (ep *entityPrototypes) buildPrototype(name string, blocks []*EntityBlock) (*entities.Entity, error) {
	e := entities.NewEntity(nil)

	processedComponents, processedRules, err := ep.expandBlocks(name, blocks)
	if err != nil {
		return nil, fmt.Errorf("could not build prototype: %w", err)
	}

	for _, c := range processedComponents {
		e.Add(c)
	}

	if len(processedRules) > 0 {
		// add eventful to component if it doesn't already have it.
		eventful, ok := entities.GetComponent[*components.Eventful](e)
		if !ok {
			eventful = &components.Eventful{Rules: []*entities.Rule{}}
			e.Add(eventful)
		}

		for _, r := range processedRules {
			eventful.AddRule(r)
		}
	}

	for _, block := range blocks {
		if block.Component == nil {
			continue
		}

		for _, f := range block.Component.Fields {
			if f.Key == "children" {
				if ep.childrenPlan[name] == nil {
					ep.childrenPlan[name] = make(map[entities.ComponentType][]string)
				}

				// populate pending children map
				componentType, err := entities.ParseComponentType(block.Component.Name)
				if err != nil {
					return nil, fmt.Errorf("could not build prototype '%s': %w", name, err)
				}

				ep.childrenPlan[name][componentType] =
					append(ep.childrenPlan[name][componentType], f.Value.Strings...)
			}
		}
	}

	return e, nil
}

// recursively expand traits in entities
func (ep *entityPrototypes) expandBlocks(name string, blocks []*EntityBlock) ([]entities.Component, []*entities.Rule, error) {
	if _, ok := ep.visiting[name]; ok {
		return nil, nil, fmt.Errorf("cycle detected at %q", name)
	}
	ep.visiting[name] = struct{}{}
	defer func() { delete(ep.visiting, name) }()

	components := make([]entities.Component, 0, len(blocks))
	rules := make([]*entities.Rule, 0, len(blocks))

	for _, block := range blocks {
		// process rules
		if block.Rule != nil {
			rule, err := buildRule(block.Rule)
			if err != nil {
				return nil, nil, fmt.Errorf("could not process rule %s: %w", block.Rule.Command, err)
			}
			rules = append(rules, rule)
		} else if block.Component != nil {
			// process component into prototype without children
			comp, err := processComponentPrototype(block.Component)
			if err != nil {
				return nil, nil, fmt.Errorf("could not process component %s: %w", block.Component.Name, err)
			}
			components = append(components, comp)
		} else if block.Trait != nil {
			traitC, traitR, err := ep.expandBlocks(block.Trait.Name, ep.traitsByName[block.Trait.Name].Blocks)
			if err != nil {
				return nil, nil, fmt.Errorf("could not process trait '%s': %w", block.Trait.Name, err)
			}

			components = append(components, traitC...)
			rules = append(rules, traitR...)
		}
	}

	return components, rules, nil
}

// loop through prototypes and instantiate them into a map of entities by name
func (ep *entityPrototypes) instantiatePrototypes() (map[string]*entities.Entity, error) {
	out := make(map[string]*entities.Entity, len(ep.prototypesByName))
	for name := range ep.prototypesByName {
		entity, err := ep.instantiate(name, nil)
		if err != nil {
			return nil, fmt.Errorf("could not instantiate '%s': %w", name, err)
		}
		out[name] = entity
	}
	return out, nil
}

// recursively instantiate a named prototype and wire up children for all child-holding components.
func (ep *entityPrototypes) instantiate(name string, parent entities.ComponentWithChildren) (*entities.Entity, error) {
	be, ok := ep.prototypesByName[name]
	if !ok {
		return nil, fmt.Errorf("unknown prototype %q", name)
	}
	if _, ok := ep.visiting[name]; ok {
		return nil, fmt.Errorf("cycle detected at %q", name)
	}
	ep.visiting[name] = struct{}{}
	defer func() { delete(ep.visiting, name) }()

	inst := be.ent.Copy(parent)

	// for each child-capable component on the entity, look up its pending child names from the prototype’s sidecar and attach recursively.
	if rm, ok := entities.GetComponent[*components.Room](inst); ok {
		slot := ep.childrenPlan[name][entities.ComponentRoom]
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
