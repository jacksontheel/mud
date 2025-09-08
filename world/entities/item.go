package entities

type Item struct {
	Id          string
	Description string
	Aliases     []string
}

var _ Entity = (*Item)(nil)

func (i *Item) GetDescription() string {
	return i.Description
}

func (i *Item) GetAliases() []string {
	return i.Aliases
}
