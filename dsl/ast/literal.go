package ast

type Literal struct {
	String  *string  `  @String`
	Bool    *string  `| @Bool`
	Strings []string `| "[" @String { "," @String } "]"`
	Pairs   []KV     `| "{" @@ { "," @@ } "}"`
}

func (l *Literal) Parse() any {
	if l.String != nil {
		return *l.String
	} else if l.Bool != nil {
		return *l.Bool == "true"
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
