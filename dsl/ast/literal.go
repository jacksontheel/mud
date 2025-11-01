package ast

import "strings"

type Literal struct {
	String  *string  `parser:"  @String"`
	Bool    *string  `parser:"| @Bool"`
	Number  *int     `parser:"| @Int"`
	Strings []string `parser:"| '[' @String { ',' @String } ']'"`
	Pairs   []KV     `parser:"| '{' @@ { ',' @@ } '}'"`
}

func (l *Literal) Parse() any {
	if l.String != nil {
		return *l.String
	} else if l.Bool != nil {
		return *l.Bool == "true"
	} else if l.Number != nil {
		return *l.Number
	} else if len(l.Strings) > 0 {
		return l.Strings
	} else {
		return l.Pairs
	}
}

func (l *Literal) AsMap() map[string]string {
	if len(l.Pairs) == 0 {
		return nil
	}
	m := make(map[string]string, len(l.Pairs))
	for _, kv := range l.Pairs {
		m[kv.Key] = kv.Value
	}
	return m
}

func (l *Literal) UnquotedString() string {
	if l.String == nil {
		return ""
	}
	return strings.Trim(*l.String, `"`)
}
