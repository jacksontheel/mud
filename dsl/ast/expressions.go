package ast

// Expressions are adapted from a Participle example (https://github.com/alecthomas/participle/blob/master/_examples/expr2/main.go)

type Expression struct {
	Equality *Equality `@@`
}

type Equality struct {
	Comparison *Comparison `@@`
	Op         string      `( @( "!" "=" | "==" )`
	Next       *Equality   `  @@ )*`
}

type Comparison struct {
	Addition *Addition   `@@`
	Op       string      `( @( ">" | ">=" | "<" | "<=" )`
	Next     *Comparison `  @@ )*`
}

type Addition struct {
	Multiplication *Multiplication `@@`
	Op             string          `( @( "-" | "+" )`
	Next           *Addition       `  @@ )*`
}

type Multiplication struct {
	Unary *Unary          `@@`
	Op    string          `( @( "/" | "*" )`
	Next  *Multiplication `  @@ )*`
}

type Unary struct {
	Op      string   `  ( @( "!" | "-" )`
	Unary   *Unary   `    @@ )`
	Primary *Primary `| @@`
}

type Primary struct {
	Number        *int        `  @Int`
	String        *string     `| @String`
	Bool          *bool       `| @( "true" | "false" )`
	Field         *Field      `| @@`
	SubExpression *Expression `| "(" @@ ")" `
	Nil           bool        `| @"nil"`
}

type Field struct {
	Role string `@Ident`
	Name string `"." @Ident`
}
