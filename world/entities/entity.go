package entities

import (
	"fmt"
	"reflect"
	"sync"
)

type Entity struct {
	mu         sync.RWMutex
	parent     ComponentWithChildren
	components map[reflect.Type]Component
}

func NewEntity(parent ComponentWithChildren) *Entity {
	return &Entity{
		parent:     parent,
		components: map[reflect.Type]Component{},
	}
}

func (e *Entity) Copy(parent ComponentWithChildren) *Entity {
	newEntity := NewEntity(parent)
	for _, c := range e.components {
		newEntity.Add(c.Copy())
	}
	return newEntity
}

func (e *Entity) GetParent() ComponentWithChildren {
	return e.parent
}

func (e *Entity) SetParent(parent ComponentWithChildren) {
	e.parent = parent
}

func (e *Entity) Add(c Component) *Entity {
	e.mu.Lock()
	e.components[reflect.TypeOf(c)] = c
	e.mu.Unlock()
	return e
}

func (e *Entity) GetComponentWithChildren(ct ComponentType) (ComponentWithChildren, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, c := range e.components {
		id := c.Id()
		if cwc, ok := any(c).(ComponentWithChildren); ok {
			if id == ct {
				return cwc, true
			}
		}
	}

	return nil, false
}

func (e *Entity) RequireComponentWithChildren(ct ComponentType) (ComponentWithChildren, error) {
	c, ok := e.GetComponentWithChildren(ct)

	if !ok {
		return nil, fmt.Errorf("entity does not have component with children %s", ct.String())
	}

	return c, nil
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

func RequireComponent[C Component](e *Entity) (C, error) {
	c, ok := GetComponent[C](e)

	if ok {
		return c, nil
	}

	var zero C
	return zero, fmt.Errorf("entity does not have component %s", zero.Id().String())
}
