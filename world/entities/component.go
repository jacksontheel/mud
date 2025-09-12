package entities

// all items have a description
type Descriptioned interface {
	Description() string
}

type CDescriptioned struct {
	Text string `json:"text"`
}

var _ Descriptioned = &CDescriptioned{}

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

var _ Aliased = &CAliased{}

func (c *CAliased) Aliases() []string {
	return c.Text
}

// most items have a tag
type Tagged interface {
	Tags() []string
}

type CTagged struct {
	Text []string `json:"text"`
}

var _ Tagged = &CTagged{}

func (c *CTagged) Tags() []string {
	return c.Text
}
