package ast

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var DslLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Bool", Pattern: `\b(?:true|false)\b`},
	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "AtIdent", Pattern: `@[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Tag", Pattern: `#[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "String", Pattern: `"([^"\\]|\\.)*"`},
	{Name: "Dot", Pattern: `\.`},
	{Name: "LBrack", Pattern: `\[`},
	{Name: "RBrack", Pattern: `\]`},
	{Name: "Colon", Pattern: `:`},
	{Name: "Comma", Pattern: `,`},
	{Name: "LBrace", Pattern: `{`},
	{Name: "RBrace", Pattern: `}`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
})

type DSL struct {
	Declarations []*TopLevel `@@*`
}

type TopLevel struct {
	Entity *EntityDef `"entity" @@`
	Trait  *TraitDef  `| "trait" @@`
}

type EntityDef struct {
	Name   string         `@Ident`
	Blocks []*EntityBlock `"{" { @@ } "}"`
}

type TraitDef struct {
	Name   string         `@Ident`
	Blocks []*EntityBlock `"{" { @@ } "}"`
}

type EntityBlock struct {
	Component *ComponentDef        `  "component" @@ `
	Trait     *TraitInheritanceDef `| "trait" @@`
	Reaction  *ReactionDef         `| "react" @@ `
	Field     *FieldDef            `| @@`
}

type ComponentDef struct {
	Name   string      `@Ident`
	Fields []*FieldDef `"{" { @@ } "}"`
}

type TraitInheritanceDef struct {
	Name string `@Ident`
}

type ReactionDef struct {
	Command string     `@Ident`
	Rules   []*RuleDef `"{" { @@ } "}"`
}

type RuleDef struct {
	When *WhenBlock `[ "when" @@ ]`
	Then *ThenBlock `"then" @@`
}

type WhenBlock struct {
	Conds []*ConditionDef `"{" { @@ } "}"`
}

type ThenBlock struct {
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
