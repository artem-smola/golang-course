package usecase

import (
	"context"

	"github.com/artem-smola/golang-course/task4/api/internal/domain"
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

func (u *Ping) Execute(ctx context.Context) domain.OverallPingResponse {
	processorStatus := u.processorPinger.Ping(ctx)
	subscriberStatus := u.subscriberPinger.Ping(ctx)

	status := domain.OverallStatusOK
	if processorStatus == domain.PingStatusDown || subscriberStatus == domain.PingStatusDown {
		status = domain.OverallStatusDegraded
	}
	response := domain.OverallPingResponse{
		Status: status,
		Services: []domain.ServicePingResponse{
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

	return response
}
