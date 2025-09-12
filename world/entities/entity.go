package entities

import (
	"strings"
	"sync"
)

type Entity struct {
	mu         sync.RWMutex
	components []any
}

func NewEntity() *Entity {
	return &Entity{components: []any{}}
}

func (e *Entity) Add(c any) *Entity {
	e.mu.Lock()
	e.components = append(e.components, c)
	e.mu.Unlock()
	return e
}

func Find[T any](e *Entity) (T, bool) {
	var zero T
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, c := range e.components {
		if v, ok := c.(T); ok {
			return v, true
		}
	}
	return zero, false
}

// util funcs for entities with common components
func (e *Entity) getAliases() []string {
	if a, ok := Find[Aliased](e); ok {
		return a.Aliases()
	}
	return nil
}

func (e *Entity) getDescription() string {
	if d, ok := Find[Descriptioned](e); ok {
		return strings.TrimSpace(d.Description())
	}
	return ""
}

func (e *Entity) getTags() []string {
	if t, ok := Find[Tagged](e); ok {
		return t.Tags()
	}
	return nil
}
