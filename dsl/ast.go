package dsl

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var dslLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "AtIdent", Pattern: `@[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Tag", Pattern: `#[a-zA-Z_][a-zA-Z0-9_]*`},
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
	Name   string         `"entity" @Ident`
	Blocks []*EntityBlock `"{" { @@ } "}"`
}

type EntityBlock struct {
	Component *ComponentDef `  "has" @@ `
	Rule      *RuleDef      `| "when" @@ `
}

type ComponentDef struct {
	Name   string      `@Ident`
	Fields []*FieldDef `"{" { @@ } "}"`
}

type RuleDef struct {
	Command string       `@Ident`
	By      *string      `[ "by" @Tag ]`
	With    *string      `[ "with" @Tag ]`
	On      *string      `[ "on" @Tag ]`
	Actions []*ActionDef `"{" { @@ } "}"`
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

type ActionDef struct {
	Say *SayAction `  @@`
}

type SayAction struct {
	Value string `"say" @String`
}
