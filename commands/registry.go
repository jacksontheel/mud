package commands

import (
	"fmt"
	"log"

	"example.com/mud/commands/cmd"
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

	all = append(all, cmd.MovePatterns...)
	all = append(all, cmd.TakePatterns...)
	all = append(all, cmd.LookPatterns...)
	all = append(all, cmd.AttackPatterns...)
	all = append(all, cmd.KissPatterns...)
	all = append(all, cmd.InventoryPatterns...)
	all = append(all, cmd.SayPatterns...)
	all = append(all, cmd.WhisperPatterns...)

	return all
}

func AllAliases() map[string]string {

	var merged, err = MergeAliasMaps(
		cmd.MoveAliases,
		cmd.TakeAliases,
		cmd.LookAliases,
		cmd.AttackAliases,
		cmd.KissAliases,
		cmd.InventoryAliases,
		cmd.SayAliases,
		cmd.WhisperAliases,
	)

	if err != nil {
		log.Fatalf("Failed to merge alias maps: %v", err)
	}

	return merged
}

func DirectionAliases() map[string]string {
	return cmd.MoveDirectionAliases
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
