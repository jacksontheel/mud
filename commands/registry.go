package commands

import (
	"fmt"
	"log"

	"example.com/mud/models"
)

const (
	CommandMove      = "move"
	CommandTake      = "take"
	CommandLook      = "look"
	CommandAttack    = "attack"
	CommandKiss      = "kiss"
	CommandInventory = "inventory"
	CommandSay       = "say"
	CommandWhisper   = "whisper"
)

func AllPatterns() []models.Pattern {
	var all []models.Pattern

	all = append(all, movePatterns...)
	all = append(all, takePatterns...)
	all = append(all, lookPatterns...)
	all = append(all, attackPatterns...)
	all = append(all, kissPatterns...)
	all = append(all, inventoryPatterns...)
	all = append(all, sayPatterns...)
	all = append(all, whisperPatterns...)

	return all
}

func AllAliases() map[string]string {

	var merged, err = MergeAliasMaps(
		moveAliases,
		takeAliases,
		lookAliases,
		attackAliases,
		kissAliases,
		inventoryAliases,
		sayAliases,
		whisperAliases,
	)

	if err != nil {
		log.Fatalf("Failed to merge alias maps: %v", err)
	}

	return merged
}

func DirectionAliases() map[string]string {
	return moveDirectionAliases
}

func MergeAliasMaps(aliasMaps ...map[string]string) (map[string]string, error) {
	combinedAliasMap := make(map[string]string)
	for _, m := range aliasMaps {
		for k, v := range m {
			if existing, ok := combinedAliasMap[k]; ok && existing != v {
				return nil, fmt.Errorf("alias conflict: %q maps to both %q and %q", k, existing, v)
			}
			combinedAliasMap[k] = v
		}
	}
	return combinedAliasMap, nil
}
