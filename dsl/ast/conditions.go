package ast

type ConditionDef struct {
	Not                      *NotCondition             `  @@`
	HasTag                   *HasTagCondition          `| @@`
	IsPresent                *IsPresentCondition       `| @@`
	EventRolesEqualCondition *EventRolesEqualCondition `| @@`
	HasChildCondition        *HasChildCondition        `| @@`
}

type NotCondition struct {
	Cond *ConditionDef `"not" @@`
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

type HasChildCondition struct {
	ParentRole string `@Ident`
	Component  string `@Ident`
	ChildRole  string `"has" "child" @Ident`
}
