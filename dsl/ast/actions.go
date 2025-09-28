package ast

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
	EntityId  string `@String`
	Target    string `"to" @Ident `
	Component string `"." @Ident`
}

type MoveAction struct {
	RoleOrigin      string ` @Ident`
	RoleDestination string ` "to" @Ident`
	Component       string `"." @Ident`
}
