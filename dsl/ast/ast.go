package ast

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var DslLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "String", Pattern: `"([^"\\]|\\.)*"`},
	{Name: "Bool", Pattern: `\b(?:true|false)\b`},
	{Name: "Int", Pattern: `[-+]?[0-9]+`},

	{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "AtIdent", Pattern: `@[a-zA-Z_][a-zA-Z0-9_]*`},
	{Name: "Tag", Pattern: `#[a-zA-Z_][a-zA-Z0-9_]*`},

	{Name: "Op", Pattern: `==|!=|<=|>=|\|\||&&`},
	{Name: "Punct", Pattern: `[\[\](){}.,=<>+\-*/!]|:`},

	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
})

type DSL struct {
	Declarations []*TopLevel `@@*`
}

type TopLevel struct {
	Entity  *EntityDef  `"entity" @@`
	Trait   *TraitDef   `| "trait" @@`
	Command *CommandDef `| "command" @@`
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
	Name   string      `@Ident`
	Fields []*FieldDef `("{" { @@ } "}")?`
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

type CommandDef struct {
	Name   string          `@Ident`
	Blocks []*CommandBlock `"{" { @@ } "}"`
}

type CommandBlock struct {
	Field                *FieldDef             `  @@`
	CommandDefinitionDef *CommandDefinitionDef `| @@`
}

type CommandDefinitionDef struct {
	Fields []*FieldDef `"pattern" "{" { @@ } "}"`
}
