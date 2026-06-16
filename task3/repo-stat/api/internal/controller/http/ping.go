package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := ping.Execute(r.Context())
		services := make([]dto.ServicePingResponse, 0, len(result.Services))
		for _, service := range result.Services {
			services = append(services, dto.ServicePingResponse{
				Name:   string(service.Name),
				Status: string(service.Status),
			})
		}

		response := dto.PingResponse{
			Status:   string(result.Status),
			Services: services,
		}

		w.Header().Set("Content-Type", "application/json")
		if result.Status == domain.OverallStatusDegraded {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
