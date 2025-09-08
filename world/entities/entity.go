package entities

type Entity interface {
	GetAliases() []string
	GetDescription() string
}

type IsEdible interface {
	GetEatDescription() string
}
