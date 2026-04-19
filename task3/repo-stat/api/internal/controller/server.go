package controller

import (
	"errors"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RESTServer struct {
	log             *slog.Logger
	router          *gin.Engine
	pingUseCase     *usecase.Ping
	repoInfoUseCase *usecase.GetRepoInfo
}

func NewRESTServer(log *slog.Logger, pingUseCase *usecase.Ping, repoInfoUseCase *usecase.GetRepoInfo) *RESTServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	server := &RESTServer{
		log:             log,
		router:          router,
		pingUseCase:     pingUseCase,
		repoInfoUseCase: repoInfoUseCase,
	}

	server.registerRoutes()
	return server
}

func (s *RESTServer) registerRoutes() {
	s.router.GET("/api/ping", s.Ping)
	s.router.GET("/api/repositories/info", s.GetRepoInfo)
}

func (s *RESTServer) Handler() http.Handler {
	return s.router
}

// Ping godoc
// @Summary Service health check
// @Description Checks the availability of processor and subscriber services.
// @Tags health
// @Produce json
// @Success 200 {object} dto.PingResponse
// @Failure 503 {object} dto.PingResponse
// @Router /api/ping [get]
func (s *RESTServer) Ping(c *gin.Context) {
	result := s.pingUseCase.Execute(c.Request.Context())
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

	if response.Status == "degraded" {
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetRepoInfo godoc
// @Summary Get GitHub repository info
// @Description Returns basic information for a GitHub repository by URL.
// @Tags repositories
// @Produce json
// @Param url query string true "GitHub repository URL"
// @Success 200 {object} dto.RepoInfo
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/repositories/info [get]
func (s *RESTServer) GetRepoInfo(c *gin.Context) {
	repoURL := c.Query("url")
	repoInfo, err := s.repoInfoUseCase.Execute(repoURL)
	if err != nil {
		if errors.Is(err, usecase.ErrEmptyRepositoryURL) || errors.Is(err, usecase.ErrInvalidRepositoryURL) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		s.log.Error("failed to get repository info", "error", err)
		switch status.Code(err) {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid repository request"})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "repository not found"})
		case codes.DeadlineExceeded, codes.Unavailable:
			c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{Error: "repository service unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get repository info"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.RepoInfo{
		FullName:    repoInfo.Name,
		Description: repoInfo.Description,
		Stars:       int64(repoInfo.StarsCount),
		Forks:       int64(repoInfo.ForksCount),
		CreatedAt:   repoInfo.CreatedAt,
	})
}
