package entities

type Food struct {
	Id             string
	Description    string
	Aliases        []string
	EatDescription string
}

// compile-time checks
var _ Entity = (*Food)(nil)
var _ IsEdible = (*Food)(nil)

func (f *Food) GetDescription() string {
	return f.Description
}

func (f *Food) GetAliases() []string {
	return f.Aliases
}

func (f *Food) GetEatDescription() string {
	return f.EatDescription
}
