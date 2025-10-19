package ast

type ConditionDef struct {
	Or *OrChain `@@`
}

type OrChain struct {
	First *CondAtom     `@@`
	Rest  []*OrChainRhs `( "or" @@ )*`
}

type OrChainRhs struct {
	Next *CondAtom `@@`
}

type CondAtom struct {
	Paren      *ConditionDef             `  "(" @@ ")"`
	Not        *NotCondition             `| @@`
	Expr       *ExprCondition            `| @@`
	HasTag     *HasTagCondition          `| @@`
	IsPresent  *IsPresentCondition       `| @@`
	RolesEqual *EventRolesEqualCondition `| @@`
	FieldEq    *FieldEqualsCondition     `| @@`
	HasChild   *HasChildCondition        `| @@`
	MsgHas     *MessageContains          `| @@`
}

type NotCondition struct {
	Cond *ConditionDef `"not" @@`
}

type ExprCondition struct {
	Expr *Expression `"expr" "{" @@ "}"`
}

type HasTagCondition struct {
	Target string `@Ident`
	Tag    string `"has" "tag" @String`
}

type IsPresentCondition struct {
	Role string `@Ident "exists"`
}

type EventRolesEqualCondition struct {
	Role1 string `@Ident`
	Role2 string `"is" @Ident`
}

type FieldEqualsCondition struct {
	Role  string  `@Ident`
	Field string  `"." @Ident`
	Value Literal `"is" @@`
}

type HasChildCondition struct {
	ChildRole  string `@Ident`
	ParentRole string `"in" @Ident`
	Component  string `"." @Ident`
}

type MessageContains struct {
	Message string `"message" "contains" @String`
}
