package commands

import (
	"fmt"

	"example.com/mud/models"
)

var Commands = map[string]struct{}{}

var DirectionAliases = map[string]string{
	"n":     models.DirectionNorth,
	"north": models.DirectionNorth,

	"s":     models.DirectionSouth,
	"south": models.DirectionSouth,

	"e":    models.DirectionEast,
	"east": models.DirectionEast,

	"w":    models.DirectionWest,
	"west": models.DirectionWest,

	"u":  models.DirectionUp,
	"up": models.DirectionUp,

	"d":    models.DirectionDown,
	"down": models.DirectionDown,
}

var VerbAliases = map[string]string{}

var Patterns = []models.Pattern{}

func init() {
	addCommandDefinition(attackCommand)
	addCommandDefinition(inventoryCommand)
	addCommandDefinition(kissCommand)
	addCommandDefinition(lookCommand)
	addCommandDefinition(moveCommand)
	addCommandDefinition(sayCommand)
	addCommandDefinition(takeCommand)
	addCommandDefinition(whisperCommand)
}

func addCommandDefinition(cd CommandDefinition) error {
	if len(cd.Aliases) == 0 {
		return fmt.Errorf("command definition must have at least one alias")
	}

	// add command to commands
	Commands[cd.Aliases[0]] = struct{}{}

	// add aliases to verbAliases
	canonicalName := cd.Aliases[0]
	for _, a := range cd.Aliases {
		VerbAliases[a] = canonicalName
	}

	// add patterns to patterns
	Patterns = append(Patterns, cd.Pattern...)

	return nil
}

// func MergeAliasMaps(aliasMaps ...map[string]string) (map[string]string, error) {
// 	combinedAliasMap := make(map[string]string)
// 	for _, m := range aliasMaps {
// 		for k, v := range m {
// 			if existing, ok := combinedAliasMap[k]; ok && existing != v {
// 				return nil, fmt.Errorf("alias conflict: %q maps to both %q and %q", k, existing, v)
// 			}
// 			combinedAliasMap[k] = v
// 		}
// 	}
// 	return combinedAliasMap, nil
// }
