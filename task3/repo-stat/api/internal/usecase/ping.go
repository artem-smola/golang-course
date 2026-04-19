package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}

type Ping struct {
	subscriberPinger Pinger
	processorPinger  Pinger
}

func NewPing(subscriberPinger Pinger, processorPinger Pinger) *Ping {
	return &Ping{
		subscriberPinger: subscriberPinger,
		processorPinger:  processorPinger,
	}
}

func (u *Ping) Execute(ctx context.Context) domain.PingResponse {
	processorStatus := u.processorPinger.Ping(ctx)
	subscriberStatus := u.subscriberPinger.Ping(ctx)

	response := domain.PingResponse{
		Status: domain.OverallStatusOK,
		Services: []domain.ServiceResponse{
			{
				Name:   domain.ServiceNameProcessor,
				Status: processorStatus,
			},
			{
				Name:   domain.ServiceNameSubscriber,
				Status: subscriberStatus,
			},
		},
	}

	if processorStatus == domain.PingStatusDown || subscriberStatus == domain.PingStatusDown {
		response.Status = domain.OverallStatusDegraded
	}

	return response
}
