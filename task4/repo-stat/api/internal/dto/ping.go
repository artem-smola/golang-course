package dto

type ServicePingResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponse struct {
	Status   string                `json:"status"`
	Services []ServicePingResponse `json:"services"`
}
