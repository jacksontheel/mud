package loading

import (
	"encoding/json"

	"example.com/mud/world/entities"
)

type rawRoom struct {
	Id          string            `json:"id"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
	RawItems    []rawItem         `json:"items"`
}

type rawItem struct {
	Id         string            `json:"id"`
	Components []json.RawMessage `json:"components"`
}

type rawComponent struct {
	Type string `json:"type"`
}

type rawEventful struct {
	Rules []rawRule `json:"rules"`
}

type rawRule struct {
	When entities.When     `json:"when"`
	Then []json.RawMessage `json:"then"`
}

type rawAction struct {
	Type string `json:"type"`
}
