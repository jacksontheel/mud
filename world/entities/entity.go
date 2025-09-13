package entities

import (
	"reflect"
	"strings"
	"sync"
)

type Entity struct {
	mu         sync.RWMutex
	components map[reflect.Type]any
}

func NewEntity() *Entity {
	return &Entity{components: map[reflect.Type]any{}}
}

func (e *Entity) Add(c any) *Entity {
	e.mu.Lock()
	e.components[reflect.TypeOf(c)] = c
	e.mu.Unlock()
	return e
}

func GetComponent[C Component](e *Entity) (C, bool) {
	var zero C
	e.mu.RLock()
	defer e.mu.RUnlock()
	v, ok := e.components[reflect.TypeOf((*C)(nil)).Elem()]
	if !ok {
		return zero, false
	}
	return v.(C), true
}

// util funcs for entities with common components
func (e *Entity) getName() string {
	if n, ok := GetComponent[*Named](e); ok {
		return strings.TrimSpace(n.Name)
	}
	return ""
}

func (e *Entity) getAliases() []string {
	if a, ok := GetComponent[*Aliased](e); ok {
		return a.Aliases
	}
	return nil
}

func (e *Entity) getDescription() string {
	if d, ok := GetComponent[*Descriptioned](e); ok {
		return strings.TrimSpace(d.Description)
	}
	return ""
}

func (e *Entity) getTags() []string {
	if t, ok := GetComponent[*Tagged](e); ok {
		return t.Tags
	}
	return nil
}
