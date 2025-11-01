package ast

// Expressions are adapted from a Participle example (https://github.com/alecthomas/participle/blob/master/_examples/expr2/main.go)

type Expression struct {
	Equality *Equality `parser:"@@"`
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
}

type Field struct {
	Role string `parser:"@Ident"`
	Name string `parser:"( '.' @Ident )?"`
}
