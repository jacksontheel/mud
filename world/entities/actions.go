package entities

type Action interface {
	Id() string
}

type Say interface {
	Say() string
}

type ASay struct {
	Text string `json:"text"`
}

func (a *ASay) Id() string {
	return "say"
}

func (a *ASay) Say() string {
	return a.Text
}
