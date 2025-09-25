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
	Rule      *RuleDef             `| "when" @@ `
	Field     *FieldDef            `| @@`
}

type ComponentDef struct {
	Name   string      `@Ident`
	Fields []*FieldDef `"{" { @@ } "}"`
}

type TraitInheritanceDef struct {
	Name string `@Ident`
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
	Print   *PrintAction   `  "print" @@`
	Publish *PublishAction `| "publish" @@`
	Copy    *CopyAction    `| "copy" @@`
	Move    *MoveAction    `| "move" @@`
}

type PrintAction struct {
	Target string `@Ident`
	Value  string ` @String`
}

type PublishAction struct {
	Value string ` @String`
}

type CopyAction struct {
	EntityId  string ` @String`
	Target    string ` "to" @Ident `
	Component string ` @Ident`
}

type MoveAction struct {
	RoleOrigin      string ` @Ident`
	RoleDestination string ` "to" @Ident `
	Component       string ` @Ident`
}
