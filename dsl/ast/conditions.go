package ast

type ConditionDef struct {
	Not                      *NotCondition             `  @@`
	HasTag                   *HasTagCondition          `| @@`
	IsPresent                *IsPresentCondition       `| @@`
	EventRolesEqualCondition *EventRolesEqualCondition `| @@`
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
