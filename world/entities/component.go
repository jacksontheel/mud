package entities

type Component interface {
	Id() string
}

type Named struct {
	Name string `json:"text"`
}

var _ Component = &Named{}

func (c *Named) Id() string {
	return "named"
}

type Descriptioned struct {
	Description string `json:"text"`
}

var _ Component = &Descriptioned{}

func (c *Descriptioned) Id() string {
	return "descriptioned"
}

type Aliased struct {
	Aliases []string `json:"text"`
}

var _ Component = &Aliased{}

func (c *Aliased) Id() string {
	return "aliased"
}

type Tagged struct {
	Tags []string `json:"text"`
}

var _ Component = &Tagged{}

func (c *Tagged) Id() string {
	return "tagged"
}
