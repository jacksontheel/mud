package loading

import "encoding/json"

type rawRoom struct {
	Id          string            `json:"id"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
	Entities    []json.RawMessage `json:"items"` // defer decoding
}

type entityEnvelope struct {
	Type string `json:"type"`
}
