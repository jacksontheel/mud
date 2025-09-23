package commands

import (
	"fmt"
	"log"
)

func AllPatterns() []Pattern {
	var all []Pattern

	all = append(all, movePatterns...)
	all = append(all, takePatterns...)
	all = append(all, lookPatterns...)

	return all
}

func AllAliases() map[string]string {

	var merged, err = MergeAliasMaps(
		moveAliases,
		takeAliases,
		lookAliases,
	)

	if err != nil {
		log.Fatalf("Failed to merge alias maps: %v", err)
	}

	return merged
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
