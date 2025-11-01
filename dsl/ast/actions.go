package ast

type ActionDef struct {
	Print                *PrintAction          `parser:"  'print' @@"`
	Publish              *PublishAction        `parser:"| 'publish' @@"`
	Copy                 *CopyAction           `parser:"| 'copy' @@"`
	Move                 *MoveAction           `parser:"| 'move' @@"`
	SetField             *SetFieldAction       `parser:"| 'set' @@"`
	DestroyAction        *DestroyAction        `parser:"| 'destroy' @@"`
	RevealChildrenAction *RevealChildrenAction `parser:"| @@"`
}

type PrintAction struct {
	Target string `parser:"@Ident"`
	Value  string `parser:"@String"`
}

type PublishAction struct {
	Value string `parser:"@String"`
}

type CopyAction struct {
	EntityId  string `parser:"@String"`
	Target    string `parser:"'to' @Ident"`
	Component string `parser:"'.' @Ident"`
}

type MoveAction struct {
	RoleObject      string `parser:"@Ident"`
	RoleDestination string `parser:"'to' @Ident"`
	Component       string `parser:"'.' @Ident"`
}

type SetFieldAction struct {
	Role  string     `parser:"@Ident"`
	Field string     `parser:"'.' @Ident"`
	Expr  Expression `parser:"'to' @@"`
}

type RevealChildrenAction struct {
	Set       string `parser:"@('reveal' | 'hide')"`
	Role      string `parser:"@Ident"`
	Component string `parser:"'.' @Ident"`
}

type DestroyAction struct {
	Role string `parser:"@Ident"`
}
