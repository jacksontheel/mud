package entities

// all items have a description
type Descriptioned interface {
	Description() string
}

type CDescriptioned struct {
	Text string `json:"text"`
}

func (c *CDescriptioned) Description() string {
	return c.Text
}

// all items have an alias
type Aliased interface {
	Aliases() []string
}

type CAliased struct {
	Text []string `json:"text"`
}

func (c *CAliased) Aliases() []string {
	return c.Text
}
