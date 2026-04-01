package usecase

type Ping struct{}

func NewPing() *Ping {
	return &Ping{}
}

func (u *Ping) Execute() string {
	return "up"
}
