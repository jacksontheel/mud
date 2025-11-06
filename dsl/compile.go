package dsl

import (
	"fmt"
	"strings"

	"example.com/mud/dsl/ast"
	"example.com/mud/models"
	"example.com/mud/world/entities"
	"example.com/mud/world/entities/actions"
	"example.com/mud/world/entities/components"
	"example.com/mud/world/entities/conditions"
)

type collectedDefs struct {
	entitiesById map[string]*ast.EntityDef
	traitsById   map[string]*ast.TraitDef
	commandsById map[string]*ast.CommandDef
}

type ChildrenPlan map[string]map[entities.ComponentType][]string

type LoweredEntity struct {
	name           string
	description    string
	tags           []string
	aliases        []string
	components     []entities.Component
	fields         map[string]models.Value
	rulesByCommand map[string][]*entities.Rule
}

type entityPrototype struct {
	id  string
	ent *entities.Entity
	def *ast.EntityDef
}

type entityPrototypes struct {
	prototypesById map[string]*entityPrototype
	traitsById     map[string]*ast.TraitDef
	childrenPlan   ChildrenPlan
	visiting       map[string]struct{}
}

func Compile(ast *ast.DSL) (map[string]*entities.Entity, []*models.CommandDefinition, error) {
	if ast == nil {
		return nil, nil, fmt.Errorf("nil DSL")
	}

	collectedDefs, err := collectDefs(ast.Declarations)
	if err != nil {
		return nil, nil, fmt.Errorf("could not collect top level declarations: %w", err)
	}

	prototypes, err := collectedDefs.collectPrototypes()
	if err != nil {
		return nil, nil, fmt.Errorf("could not collect prototype entities: %w", err)
	}

	entitiesById, err := prototypes.instantiatePrototypes()
	if err != nil {
		return nil, nil, fmt.Errorf("could not instantiate prototype entities: %w", err)
	}

	commands := make([]*models.CommandDefinition, 0, len(collectedDefs.commandsById))
	for _, c := range collectedDefs.commandsById {
		cd, err := buildCommandDefinition(c)
		if err != nil {
			return nil, nil, fmt.Errorf("could not instantiate command '%s': %w", c.Name, err)
		}

		commands = append(commands, cd)
	}

	return entitiesById, commands, nil
}

// collect entity, command and trait definitions
func collectDefs(decls []*ast.TopLevel) (*collectedDefs, error) {
	entitiesById := make(map[string]*ast.EntityDef, len(decls))
	commandsById := make(map[string]*ast.CommandDef, len(decls))
	traitsById := make(map[string]*ast.TraitDef, len(decls))

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
		} else if ec := declaration.Command; ec != nil {
			if _, exists := commandsById[ec.Name]; exists {
				return nil, fmt.Errorf("duplicate command %s", ec.Name)
			}

			commandsById[declaration.Command.Name] = declaration.Command
		} else {
			return nil, fmt.Errorf("declaration at top level is empty")
		}
	}

	return &collectedDefs{
		entitiesById: entitiesById,
		traitsById:   traitsById,
		commandsById: commandsById,
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
func (ep *entityPrototypes) buildPrototype(id string, blocks []*ast.EntityBlock) (*entities.Entity, error) {

	loweredEntity, err := ep.lowerEntity(id, blocks)
	if err != nil {
		return nil, fmt.Errorf("could not build prototype: %w", err)
	}

	e := entities.NewEntity(
		loweredEntity.name,
		loweredEntity.description,
		loweredEntity.aliases,
		loweredEntity.tags,
		loweredEntity.fields,
		nil,
	)

	for _, c := range loweredEntity.components {
		e.Add(c)
	}

	if len(loweredEntity.rulesByCommand) > 0 {
		// create eventful if it doesn't already exist
		eventful, ok := entities.GetComponent[*components.Eventful](e)
		if !ok {
			eventful = &components.Eventful{
				Rules: map[string][]*entities.Rule{},
			}
			e.Add(eventful)
		}

		for command, rulesByCommand := range loweredEntity.rulesByCommand {
			for _, r := range rulesByCommand {
				eventful.AddRule(command, r)
			}
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

				// get list of strings from expression
				childrenStrings, err := immediateEvalExpressionAs(f.Value, models.KindStringList)
				if err != nil {
					return nil, fmt.Errorf("could not get children list for prototype '%s': %w", id, err)
				}

				ep.childrenPlan[id][componentType] =
					append(ep.childrenPlan[id][componentType], childrenStrings.SL...)
			}
		}
	}

	return e, nil
}

// recursively expand traits in entities
func (ep *entityPrototypes) lowerEntity(id string, blocks []*ast.EntityBlock) (*LoweredEntity, error) {
	if _, ok := ep.visiting[id]; ok {
		return nil, fmt.Errorf("cycle detected at %q", id)
	}
	ep.visiting[id] = struct{}{}
	defer func() { delete(ep.visiting, id) }()

	var name string
	var description string
	var aliases []string
	var tags []string
	fields := make(map[string]models.Value)

	components := make([]entities.Component, 0, len(blocks))
	rulesByCommand := make(map[string][]*entities.Rule, len(blocks))

	for _, block := range blocks {
		if block.Reaction != nil {
			// process reaction
			rules, err := buildReaction(block.Reaction)
			if err != nil {
				return nil, err
			}
			// rules at the entity level come first
			for _, command := range block.Reaction.Commands {
				rulesByCommand[command] = append(rules, rulesByCommand[command]...)
			}
		} else if block.Component != nil {
			// process component into prototype without children
			comp, err := processComponentPrototype(block.Component)
			if err != nil {
				return nil, fmt.Errorf("could not process component %s: %w", block.Component.Name, err)
			}
			components = append(components, comp)
		} else if block.Trait != nil {
			// TODO this dereferences a nil pointer if the trait doesn't exist
			loweredTrait, err := ep.lowerEntity(block.Trait.Name, ep.traitsById[block.Trait.Name].Blocks)
			if err != nil {
				return nil, fmt.Errorf("could not process trait '%s': %w", block.Trait.Name, err)
			}

			// first write over fields that were passed into trait
			for _, f := range block.Trait.Fields {
				value, err := immediateEvalExpression(f.Value)
				if err != nil {
					return nil, fmt.Errorf("could not get process trait '%s' field '%s': %w", block.Trait.Name, f.Key, err)
				}

				// only include fields passed into trait that aren't already defined
				if _, ok := fields[f.Key]; !ok {
					fields[f.Key] = value
				}
			}

			// first add fields that were inherited from trait
			for k, tf := range loweredTrait.fields {
				// only include fields from trait that aren't already defined
				if _, ok := fields[k]; !ok {
					fields[k] = tf
				}
			}

			components = append(components, loweredTrait.components...)
			for command, traitRules := range loweredTrait.rulesByCommand {
				// rules at the trait level come second
				rulesByCommand[command] = append(rulesByCommand[command], traitRules...)
			}

		} else if block.Field != nil {
			f := block.Field
			value, err := immediateEvalExpression(block.Field.Value)
			if err != nil {
				return nil, fmt.Errorf("could not get process field '%s' for entity '%s': %w", block.Field.Key, id, err)
			}

			switch f.Key {
			case "name":
				if value.K != models.KindString {
					return nil, fmt.Errorf("name must be a string")
				}
				name = value.S
			case "description":
				if value.K != models.KindString {
					return nil, fmt.Errorf("description must be a string")
				}
				description = value.S
			case "aliases":
				if value.K != models.KindStringList {
					return nil, fmt.Errorf("aliases must be a string list")
				}
				aliases = value.SL
			case "tags":
				if value.K != models.KindStringList {
					return nil, fmt.Errorf("tags must be a string list")
				}
				tags = value.SL
			default:
				fields[f.Key] = value
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
		name:           name,
		description:    description,
		tags:           tags,
		aliases:        aliases,
		components:     components,
		fields:         fields,
		rulesByCommand: rulesByCommand,
	}, nil
}

func buildCommandDefinition(cd *ast.CommandDef) (*models.CommandDefinition, error) {
	cmd := &models.CommandDefinition{
		Name:     strings.ToLower(cd.Name),
		Aliases:  []string{},
		Patterns: []models.CommandPattern{},
	}

	for _, b := range cd.Blocks {
		if b.Field != nil {
			f := b.Field
			switch f.Key {
			case "aliases":
				value, err := immediateEvalExpressionAs(f.Value, models.KindStringList)
				if err != nil {
					return nil, fmt.Errorf("could not get value '%s' for command aliases: %w", f.Key, err)
				}
				cmd.Aliases = append(cmd.Aliases, value.SL...)
			default:
				return nil, fmt.Errorf("unknown field '%s' in command definition", f.Key)
			}
		} else if b.CommandDefinitionDef != nil {
			commandPattern, err := buildCommandPattern(b.CommandDefinitionDef)
			if err != nil {
				return nil, fmt.Errorf("could not build command pattern: %w", err)
			}

			cmd.Patterns = append(cmd.Patterns, *commandPattern)
		} else {
			return nil, fmt.Errorf("could not expand command definition block")
		}
	}

	return cmd, nil
}

func buildCommandPattern(def *ast.CommandDefinitionDef) (*models.CommandPattern, error) {
	var p = &models.CommandPattern{
		Tokens: []models.PatToken{},
	}

	for _, f := range def.Fields {
		value, err := immediateEvalExpressionAs(f.Value, models.KindString)
		if err != nil {
			return nil, fmt.Errorf("could not get value '%s' for command: %w", f.Key, err)
		}

		switch f.Key {
		case "syntax":
			p.Tokens = tokenizeCommandSyntax(value.S)
		case "noMatch":
			p.NoMatchMessage = value.S
		case "help":
			p.HelpMessage = value.S
		default:
			err := fmt.Errorf("CommandDefinitionDef Field not recognized: %s", f.Key)
			return nil, err
		}
	}
	return p, nil
}

func tokenizeCommandSyntax(s string) []models.PatToken {
	var tokens []models.PatToken
	parts := strings.Fields(s)

	for _, part := range parts[:len(parts)-1] {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			slot := strings.Trim(part, "{}")
			tokens = append(tokens, models.Slot(slot))
		} else {
			tokens = append(tokens, models.Lit(part))
		}
	}

	lastPart := parts[len(parts)-1]
	if strings.HasPrefix(lastPart, "{") && strings.HasSuffix(lastPart, "}") {
		slot := strings.Trim(lastPart, "{}")
		if strings.Contains(slot, "...") {
			tokens = append(tokens, models.SlotRest(strings.TrimSuffix(slot, "...")))
		} else {
			tokens = append(tokens, models.Slot(slot))
		}
	} else {
		tokens = append(tokens, models.Lit(lastPart))
	}

	return tokens
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
		if len(slot) > 0 && len(rm.GetChildren().GetChildren()) == 0 {
			for _, childName := range slot {
				childInst, err := ep.instantiate(childName, rm)
				if err != nil {
					return nil, err
				}
				rm.AddChild(childInst)
			}
		}
	}

	if inventory, ok := entities.GetComponent[*components.Inventory](inst); ok {
		slot := ep.childrenPlan[id][entities.ComponentInventory]
		if len(slot) > 0 && len(inventory.GetChildren().GetChildren()) == 0 {
			for _, childName := range slot {
				childInst, err := ep.instantiate(childName, inventory)
				if err != nil {
					return nil, err
				}
				inventory.AddChild(childInst)
			}
		}
	}

	if container, ok := entities.GetComponent[*components.Container](inst); ok {
		slot := ep.childrenPlan[id][entities.ComponentContainer]
		if len(slot) > 0 && len(container.GetChildren().GetChildren()) == 0 {
			for _, childName := range slot {
				childInst, err := ep.instantiate(childName, container)
				if err != nil {
					return nil, err
				}
				container.AddChild(childInst)
			}
		}
	}

	return inst, nil
}

func buildReaction(def *ast.ReactionDef) ([]*entities.Rule, error) {
	rules := make([]*entities.Rule, 0, len(def.Rules))
	for _, r := range def.Rules {
		rule, err := buildRule(r)
		if err != nil {
			return nil, fmt.Errorf("could not build reaction for %s: %w", def.Commands[0], err)
		}

		rules = append(rules, rule)
	}
	return rules, nil
}

func buildRule(def *ast.RuleDef) (*entities.Rule, error) {
	when, err := buildWhen(def.When)
	if err != nil {
		return nil, fmt.Errorf("could not build rule: %w", err)
	}

	then, err := buildThen(def.Then)
	if err != nil {
		return nil, fmt.Errorf("could not build rule: %w", err)
	}

	return &entities.Rule{
		When: when,
		Then: then,
	}, nil
}

func buildWhen(def *ast.WhenBlock) ([]entities.Condition, error) {
	if def == nil {
		return []entities.Condition{}, nil
	}

	ret := make([]entities.Condition, len(def.Conds))

	for i, cDef := range def.Conds {
		condition, err := BuildCondition(cDef)
		if err != nil {
			return nil, fmt.Errorf("build when: %w", err)
		}
		ret[i] = condition
	}

	return ret, nil
}

func BuildCondition(def *ast.ConditionDef) (entities.Condition, error) {
	if def == nil || def.Or == nil {
		return nil, fmt.Errorf("condition in when is empty")
	}

	acc, err := buildCondAtom(def.Or.First)
	if err != nil {
		return nil, err
	}

	for _, rhs := range def.Or.Rest {
		next, err := buildCondAtom(rhs.Next)
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

func buildCondAtom(atom *ast.CondAtom) (entities.Condition, error) {
	if atom == nil {
		return nil, fmt.Errorf("empty condition atom")
	}

	if atom.Paren != nil {
		return BuildCondition(atom.Paren)
	}

	if atom.Not != nil {
		inner, err := BuildCondition(atom.Not.Cond) // NotCondition wraps *ConditionDef
		if err != nil {
			return nil, fmt.Errorf("not condition: %w", err)
		}
		return &conditions.Not{Cond: inner}, nil
	}

	if atom.Expr != nil {
		expression, err := atom.Expr.Expr.Build()
		if err != nil {
			return nil, fmt.Errorf("condition expression: %w", err)
		}
		return &conditions.ExpressionTrue{Expression: expression}, nil
	}

	if atom.HasTag != nil {
		eventRole, err := entities.ParseEventRole(atom.HasTag.Target)
		if err != nil {
			return nil, fmt.Errorf("could not build has tag condition: %w", err)
		}
		return &conditions.HasTag{
			EventRole: eventRole,
			Tag:       atom.HasTag.Tag,
		}, nil
	}

	if atom.IsPresent != nil {
		eventRole, err := entities.ParseEventRole(atom.IsPresent.Role)
		if err != nil {
			return nil, fmt.Errorf("could not build is-present condition: %w", err)
		}
		return &conditions.IsPresent{EventRole: eventRole}, nil
	}

	if atom.RolesEqual != nil {
		role1, err := entities.ParseEventRole(atom.RolesEqual.Role1)
		if err != nil {
			return nil, fmt.Errorf("event roles equal condition: %w", err)
		}
		role2, err := entities.ParseEventRole(atom.RolesEqual.Role2)
		if err != nil {
			return nil, fmt.Errorf("event roles equal condition: %w", err)
		}
		return &conditions.EventRolesEqual{
			EventRole1: role1,
			EventRole2: role2,
		}, nil
	}

	if atom.HasChild != nil {
		parentRole, err := entities.ParseEventRole(atom.HasChild.ParentRole)
		if err != nil {
			return nil, fmt.Errorf("has child condition: %w", err)
		}
		component, err := entities.ParseComponentType(atom.HasChild.Component)
		if err != nil {
			return nil, fmt.Errorf("has child condition: %w", err)
		}
		childRole, err := entities.ParseEventRole(atom.HasChild.ChildRole)
		if err != nil {
			return nil, fmt.Errorf("has child condition: %w", err)
		}
		return &conditions.HasChild{
			ParentRole:    parentRole,
			ComponentType: component,
			ChildRole:     childRole,
		}, nil
	}

	if atom.MsgHas != nil {
		return &conditions.MessageContains{
			MessageRegex: strings.ToLower(atom.MsgHas.Message),
		}, nil
	}

	return nil, fmt.Errorf("unrecognized condition atom")
}

func buildThen(def *ast.ThenBlock) ([]entities.Action, error) {
	ret := make([]entities.Action, len(def.Actions))

	for i, aDef := range def.Actions {
		var newAction entities.Action

		if aDef.Print != nil {
			eventRole, err := entities.ParseEventRole(aDef.Print.Target)
			if err != nil {
				return nil, fmt.Errorf("could not build print action: %w", err)
			}

			newAction = &actions.Print{
				Text:      aDef.Print.Value,
				EventRole: eventRole,
			}
		} else if aDef.Publish != nil {
			newAction = &actions.Publish{
				Text: aDef.Publish.Value,
			}
		} else if aDef.Copy != nil {
			eventRole, err := entities.ParseEventRole(aDef.Copy.Target)
			if eventRole == entities.EventRoleUnknown {
				return nil, fmt.Errorf("could not build copy action: %w", err)
			}

			component, err := entities.ParseComponentType(aDef.Copy.Component)
			if err != nil {
				return nil, fmt.Errorf("could not build action: %w", err)
			}

			newAction = &actions.Copy{
				EntityId:      aDef.Copy.EntityId,
				EventRole:     eventRole,
				ComponentType: component,
			}
		} else if aDef.Move != nil {
			roleObject, err := entities.ParseEventRole(aDef.Move.RoleObject)
			if err != nil {
				return nil, fmt.Errorf("could not build move action for origin: %w", err)
			}

			roleDestination, err := entities.ParseEventRole(aDef.Move.RoleDestination)
			if err != nil {
				return nil, fmt.Errorf("could not build move action for destination: %w", err)
			}

			component, err := entities.ParseComponentType(aDef.Move.Component)
			if err != nil {
				return nil, fmt.Errorf("could not build action: %w", err)
			}

			newAction = &actions.Move{
				RoleObject:      roleObject,
				RoleDestination: roleDestination,
				ComponentType:   component,
			}
		} else if aDef.SetField != nil {
			role, err := entities.ParseEventRole(aDef.SetField.Role)
			if err != nil {
				return nil, fmt.Errorf("event set field action: %w", err)
			}

			expression, err := aDef.SetField.Expr.Build()
			if err != nil {
				return nil, fmt.Errorf("expression set field action: %w", err)
			}

			newAction = &actions.SetField{
				Role:       role,
				Field:      aDef.SetField.Field,
				Expression: expression,
			}
		} else if aDef.DestroyAction != nil {
			role, err := entities.ParseEventRole(aDef.DestroyAction.Role)
			if err != nil {
				return nil, fmt.Errorf("event destroy action: %w", err)
			}

			newAction = &actions.Destroy{
				Role: role,
			}
		} else if aDef.RevealChildrenAction != nil {
			role, err := entities.ParseEventRole(aDef.RevealChildrenAction.Role)
			if err != nil {
				return nil, fmt.Errorf("could not build reveal children action for role: %w", err)
			}

			component, err := entities.ParseComponentType(aDef.RevealChildrenAction.Component)
			if err != nil {
				return nil, fmt.Errorf("could not build reveal children action: %w", err)
			}

			newAction = &actions.RevealChildren{
				Role:          role,
				ComponentType: component,
				Reveal:        aDef.RevealChildrenAction.Set == "reveal",
			}

		} else {
			return nil, fmt.Errorf("action in then is empty")
		}

		ret[i] = newAction
	}

	return ret, nil
}
