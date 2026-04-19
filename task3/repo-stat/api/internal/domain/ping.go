package domain

type PingStatus string

const (
	PingStatusUp   PingStatus = "up"
	PingStatusDown PingStatus = "down"
)

type ServiceName string

const (
	ServiceNameSubscriber ServiceName = "subscriber"
	ServiceNameProcessor  ServiceName = "processor"
)

type OverallStatus string

const (
	OverallStatusOK       OverallStatus = "ok"
	OverallStatusDegraded OverallStatus = "degraded"
)

type ServiceResponse struct {
	Name   ServiceName
	Status PingStatus
}

type PingResponse struct {
	Status   OverallStatus
	Services []ServiceResponse
}
