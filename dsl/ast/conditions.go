package ast

type ConditionDef struct {
	Or *OrChain `parser:"@@"`
}

type OrChain struct {
	First *CondAtom     `parser:"@@"`
	Rest  []*OrChainRhs `parser:"( 'or' @@ )*"`
}

type OrChainRhs struct {
	Next *CondAtom `parser:"@@"`
}

type CondAtom struct {
	Paren      *ConditionDef             `parser:"  '(' @@ ')'"`
	Not        *NotCondition             `parser:"| @@"`
	Expr       *ExprCondition            `parser:"| @@"`
	HasTag     *HasTagCondition          `parser:"| @@"`
	IsPresent  *IsPresentCondition       `parser:"| @@"`
	RolesEqual *EventRolesEqualCondition `parser:"| @@"`
	HasChild   *HasChildCondition        `parser:"| @@"`
	MsgHas     *MessageContains          `parser:"| @@"`
}

type NotCondition struct {
	Cond *ConditionDef `parser:"'not' @@"`
}

type ExprCondition struct {
	Expr *Expression `parser:"'expr' '{' @@ '}'"`
}

type HasTagCondition struct {
	Target string `parser:"@Ident"`
	Tag    string `parser:"'has' 'tag' @String"`
}

type IsPresentCondition struct {
	Role string `parser:"@Ident 'exists'"`
}

type EventRolesEqualCondition struct {
	Role1 string `parser:"@Ident"`
	Role2 string `parser:"'is' @Ident"`
}

type HasChildCondition struct {
	ChildRole  string `parser:"@Ident"`
	ParentRole string `parser:"'in' @Ident"`
	Component  string `parser:"'.' @Ident"`
}

type MessageContains struct {
	Message string `parser:"'message' 'contains' @String"`
}
