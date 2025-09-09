package entities

import (
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

// First component that satisfies interface T (e.g., Edible, Throwable).
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
