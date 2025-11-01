package ast

// Expressions are adapted from a Participle example (https://github.com/alecthomas/participle/blob/master/_examples/expr2/main.go)

type Expression struct {
	Equality *Equality `parser:"  @@"`

	// pairs are only used for room exits and are not ever passed to the evaluation
	// TODO support maps as a value in the Orbis Definition language
	Pairs []KV `parser:"| '{' @@ { ',' @@ } '}'"`
}

type Equality struct {
	Comparison *Comparison `parser:"@@"`
	Op         string      `parser:"( @( '!' '=' | '==' )"`
	Next       *Equality   `parser:"  @@ )*"`
}

type Comparison struct {
	Addition *Addition   `parser:"@@"`
	Op       string      `parser:"( @( '>' | '>=' | '<' | '<=' )"`
	Next     *Comparison `parser:"  @@ )*"`
}

type Addition struct {
	Multiplication *Multiplication `parser:"@@"`
	Op             string          `parser:"( @( '-' | '+' )"`
	Next           *Addition       `parser:"  @@ )*"`
}

type Multiplication struct {
	Unary *Unary          `parser:"@@"`
	Op    string          `parser:"( @( '/' | '*' )"`
	Next  *Multiplication `parser:"  @@ )*"`
}

type Unary struct {
	Op      string   `parser:"  ( @( '!' | '-' )"`
	Unary   *Unary   `parser:"    @@ )"`
	Primary *Primary `parser:"| @@"`
}

type Primary struct {
	Number        *int        `parser:"  @Int"`
	String        *string     `parser:"| @String"`
	Bool          *string     `parser:"| @( 'true' | 'false' )"`
	Field         *Field      `parser:"| @@"`
	SubExpression *Expression `parser:"| '(' @@ ')' "`
	Nil           bool        `parser:"| @'nil'"`
	List          *List       `parser:"| @@"`
}

type List struct {
	Numbers []int    `parser:"  '[' @Int { ',' @Int } ']'"`
	Strings []string `parser:"| '[' @String { ',' @String } ']'"`
	Bools   []string `parser:"| '[' ( 'true' | 'false' ) { ',' ( 'true' | 'false' ) } ']'"`
}

type Field struct {
	Role string `parser:"@Ident"`
	Name string `parser:"( '.' @Ident )?"`
}

// In the transition away from the old literals to expressions
// maps are no longer supported generally. This is to keep maps
// working for room exits.
func (e *Expression) AsMap() map[string]string {
	if len(e.Pairs) == 0 {
		return nil
	}
	m := make(map[string]string, len(e.Pairs))
	for _, kv := range e.Pairs {
		m[kv.Key] = kv.Value
	}
	return m
}
