// loader/loader.go
package loader

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"

	"example.com/mud/loader/ir"
	"example.com/mud/world/entities"
)

type Loader struct{ L *lua.LState }

func New() *Loader       { return &Loader{L: lua.NewState()} }
func (l *Loader) Close() { l.L.Close() }

// LoadFile: Lua → IR → runtime; children are inlined so no second pass needed.
func (l *Loader) LoadFile(path string) (map[string]*entities.Entity, error) {
	if err := l.L.DoFile(path); err != nil {
		return nil, fmt.Errorf("load lua file: %w", err)
	}
	val := l.L.Get(-1)
	l.L.Pop(1)

	root, ok := val.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("file %s did not return a table", path)
	}

	rooms := make(map[string]*entities.Entity)

	// Build each top-level entity; only keep rooms in returned map
	root.ForEach(func(k, v lua.LValue) {
		id := lua.LVAsString(k)
		eir := l.buildIREntityWithID(v, id)
		rt := eir.Build()
		rooms[id] = rt
	})

	return rooms, nil
}

// buildIREntityWithID: turns one Lua table into IR.Entity (with inline children IRs)
func (l *Loader) buildIREntityWithID(val lua.LValue, thisID string) *ir.Entity {
	t, ok := val.(*lua.LTable)
	if !ok {
		return &ir.Entity{ID: thisID}
	}

	e := &ir.Entity{
		ID:          thisID,
		Name:        getString(t, "name"),
		Description: getString(t, "description"),
		Aliases:     getStringArray(t, "aliases"),
		Tags:        getStringArray(t, "tags"),
	}

	var comps ir.Components
	if compsTbl, ok := t.RawGetString("components").(*lua.LTable); ok {
		compsTbl.ForEach(func(k, v lua.LValue) {
			switch lua.LVAsString(k) {
			case "room":
				if ct, ok := v.(*lua.LTable); ok {
					rc := &ir.RoomComponent{
						Color: getString(ct, "color"),
						Icon:  getString(ct, "icon"),
						Exits: getStringMap(ct, "exits"),
					}
					// Inline-only children: array of tables
					if arr, ok := ct.RawGetString("children").(*lua.LTable); ok {
						rc.InlineChildren = l.parseInlineChildren(arr, thisID)
					}
					comps.Room = rc
				}
			}
		})
	}
	e.Components = comps
	return e
}

// parseInlineChildren expects an array of child ENTITY TABLES (no string IDs).
func (l *Loader) parseInlineChildren(arr *lua.LTable, parentID string) []*ir.Entity {
	children := make([]*ir.Entity, 0, arr.Len())
	i := 0
	arr.ForEach(func(_, child lua.LValue) {
		i++
		ct, ok := child.(*lua.LTable)
		if !ok {
			// If you want to error on strings, do it here:
			// log/collect: fmt.Printf("%s: non-table child ignored\n", parentID)
			return
		}
		// Derive/ensure child ID (local uniqueness under parent)
		id := getString(ct, "id")
		if id == "" {
			name := getString(ct, "name")
			if name == "" {
				name = "child"
			}
			id = parentID + "_" + slug(name) + "_" + itoa(i)
		}
		childIR := l.buildIREntityWithID(ct, id)
		// Safety: strip room component from inline children (no nested rooms)
		if childIR.Components.Room != nil {
			childIR.Components.Room = nil
		}
		children = append(children, childIR)
	})
	return children
}
