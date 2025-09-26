package ast

type Literal struct {
	String  *string  `  @String`
	Strings []string `| "[" @String { "," @String } "]"`
	Pairs   []KV     `| "{" @@ { "," @@ } "}"`
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
