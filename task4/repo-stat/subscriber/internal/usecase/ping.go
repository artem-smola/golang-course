package usecase

import "context"

const up string = "up"

type Ping struct{}

func NewPing() *Ping {
	return &Ping{}
}

func (u *Ping) Execute(context.Context) string {
	return up
}
