package entities

import (
	"reflect"
	"sync"
)

type Entity struct {
	mu         sync.RWMutex
	components map[reflect.Type]Component
}

func NewEntity() *Entity {
	return &Entity{components: map[reflect.Type]Component{}}
}

func (e *Entity) Copy() *Entity {
	newEntity := NewEntity()
	for _, c := range e.components {
		e.Add(c.Copy())
	}
	return newEntity
}

func (e *Entity) Add(c Component) *Entity {
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
