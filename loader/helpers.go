// loader/helpers.go
package loader

import (
	"strings"
	"unicode"

	"example.com/mud/world/entities"
	"example.com/mud/world/entities/components"
	lua "github.com/yuin/gopher-lua"
)

func getString(t *lua.LTable, key string) string {
	if s, ok := t.RawGetString(key).(lua.LString); ok {
		return string(s)
	}
	return ""
}
func getStringArray(t *lua.LTable, key string) []string {
	if v, ok := t.RawGetString(key).(*lua.LTable); ok {
		var out []string
		v.ForEach(func(_, x lua.LValue) {
			if s, ok := x.(lua.LString); ok {
				out = append(out, string(s))
			}
		})
		return out
	}
	return nil
}
func getStringMap(t *lua.LTable, key string) map[string]string {
	if v, ok := t.RawGetString(key).(*lua.LTable); ok {
		out := make(map[string]string)
		v.ForEach(func(k, x lua.LValue) { out[lua.LVAsString(k)] = lua.LVAsString(x) })
		return out
	}
	return nil
}
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
func slug(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == ' ' || r == '-' || r == '_' || r == '/':
			b.WriteByte('_')
		case r <= unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)):
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return strings.Trim(b.String(), "_")
}

// hasRoomComponent: helper for filtering return set to rooms only.
func hasRoomComponent(e *entities.Entity) bool {
	if e == nil {
		return false
	}
	_, ok := entities.GetComponent[*components.Room](e)
	return ok
}
