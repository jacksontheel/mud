package ast

type ActionDef struct {
	Print                *PrintAction          `  "print" @@`
	Publish              *PublishAction        `| "publish" @@`
	Copy                 *CopyAction           `| "copy" @@`
	Move                 *MoveAction           `| "move" @@`
	SetField             *SetFieldAction       `| "set" @@`
	DestroyAction        *DestroyAction        `| "destroy" @@`
	RevealChildrenAction *RevealChildrenAction `| @@`
}

type PrintAction struct {
	Target string `@Ident`
	Value  string ` @String`
}

type PublishAction struct {
	Value string ` @String`
}

type CopyAction struct {
	EntityId  string `@String`
	Target    string `"to" @Ident `
	Component string `"." @Ident`
}

type MoveAction struct {
	RoleOrigin      string ` @Ident`
	RoleDestination string ` "to" @Ident`
	Component       string `"." @Ident`
}

type SetFieldAction struct {
	Role  string     `@Ident`
	Field string     `"." @Ident`
	Expr  Expression `"to" @@`
}

type RevealChildrenAction struct {
	Set       string `@("reveal" | "hide")`
	Role      string `@Ident`
	Component string `"." @Ident`
}

type DestroyAction struct {
	Role string `@Ident`
}
