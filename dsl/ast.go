package dsl

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var dslLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "String", Pattern: `"([^"\\]|\\.)*"`},
	{Name: "LBrack", Pattern: `\[`},
	{Name: "RBrack", Pattern: `\]`},
	{Name: "Colon", Pattern: `:`},
	{Name: "Comma", Pattern: `,`},
	{Name: "LBrace", Pattern: `{`},
	{Name: "RBrace", Pattern: `}`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
})

type DSL struct {
	Entities []*EntityDef `@@*`
}

type EntityDef struct {
	Name       string          `"entity" @Ident`
	Components []*ComponentDef `"{" { "has" @@ } "}"`
}

type ComponentDef struct {
	Name   string      `@Ident`
	Fields []*FieldDef `"{" { @@ } "}"`
}

type FieldDef struct {
	Key   string   `@Ident "is"`
	Value *Literal `@@`
}

type KV struct {
	Key   string `@String`
	Value string `":" @String`
}

type Literal struct {
	String  *string  `  @String`
	Strings []string `| "[" @String { "," @String } "]"`
	Pairs   []KV     `| "{" @@ { "," @@ } "}"`
}
