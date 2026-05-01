package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/artem-smola/golang-course/task4/api/internal/dto"
	"github.com/artem-smola/golang-course/task4/api/internal/usecase"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	log                      *slog.Logger
	router                   *gin.Engine
	pingUsecase              *usecase.Ping
	repoInfoUsecase          *usecase.GetRepoInfoUsecase
	subscriptionUsecase      *usecase.Subscription
	subscriptionsRepoInfoUsecase *usecase.GetSubscriptionsRepoInfoUsecase
}

func NewServer(
	log *slog.Logger,
	pingUsecase *usecase.Ping,
	repoInfoUsecase *usecase.GetRepoInfoUsecase,
	subscriptionUsecase *usecase.Subscription,
	subscriptionsRepoInfoUsecase *usecase.GetSubscriptionsRepoInfoUsecase,
) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	server := &Server{
		log:                      log,
		router:                   router,
		pingUsecase:              pingUsecase,
		repoInfoUsecase:          repoInfoUsecase,
		subscriptionUsecase:      subscriptionUsecase,
		subscriptionsRepoInfoUsecase: subscriptionsRepoInfoUsecase,
	}

	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.router.GET("/api/ping", s.Ping)
	s.router.GET("/api/repositories/info", s.GetRepoInfo)
	s.router.GET("/subscriptions/info", s.GetSubscriptionsRepoInfo)
	s.router.POST("/subscriptions", s.AddSubscription)
	s.router.GET("/subscriptions", s.GetSubscriptions)
	s.router.DELETE("/subscriptions/:owner/:repo", s.DeleteSubscription)
}

func (s *Server) Handler() http.Handler {
	return s.router
}

// Ping godoc
// @Summary Service health check
// @Description Checks the availability of processor and subscriber services.
// @Tags ping
// @Produce json
// @Success 200 {object} dto.PingResponse
// @Failure 503 {object} dto.PingResponse
// @Router /api/ping [get]
func (s *Server) Ping(c *gin.Context) {
	result := s.pingUsecase.Execute(c.Request.Context())
	servicesDomain := result.Services
	servicesDto := make([]dto.ServicePingResponse, 0, len(servicesDomain))
	for _, service := range servicesDomain {
		servicesDto = append(servicesDto, dto.ServicePingResponse{
			Name:   string(service.Name),
			Status: string(service.Status),
		})
	}

	response := dto.PingResponse{
		Status:   string(result.Status),
		Services: servicesDto,
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
// @Failure 404 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/repositories/info [get]
func (s *Server) GetRepoInfo(c *gin.Context) {
	repoURL := c.Query("url")
	repoInfo, err := s.repoInfoUsecase.Execute(repoURL)
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
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		Stars:       int64(repoInfo.StarsCount),
		Forks:       int64(repoInfo.ForksCount),
		CreatedAt:   repoInfo.CreatedAt,
	})
}

// AddSubscription godoc
// @Summary Create subscription
// @Description Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.AddSubscriptionRequest true "Subscription payload"
// @Success 201 {object} dto.Subscription
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [post]
func (s *Server) AddSubscription(c *gin.Context) {
	var request dto.AddSubscriptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	subscription, err := s.subscriptionUsecase.Add(request.Owner, request.RepoName)
	if err != nil {
		if errors.Is(err, usecase.ErrEmptyOwner) || errors.Is(err, usecase.ErrEmptyRepoName) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}

		s.respondWithGRPCError(c, err, "failed to add subscription")
		return
	}

	c.JSON(http.StatusCreated, dto.Subscription{
		ID:        subscription.ID,
		Owner:     subscription.Owner,
		RepoName:  subscription.RepoName,
		CreatedAt: subscription.CreatedAt,
	})
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Description Delete subscription
// @Tags subscriptions
// @Produce json
// @Param owner path string true "Repository owner"
// @Param repo path string true "Repository name"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{owner}/{repo} [delete]
func (s *Server) DeleteSubscription(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	err := s.subscriptionUsecase.Delete(owner, repo)
	if err != nil {
		if errors.Is(err, usecase.ErrEmptyOwner) || errors.Is(err, usecase.ErrEmptyRepoName) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}

		s.respondWithGRPCError(c, err, "failed to delete subscription")
		return
	}

	c.Status(http.StatusNoContent)
}

// GetSubscriptions godoc
// @Summary List subscriptions
// @Description List subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {object} dto.SubscriptionsResponse
// @Failure 503 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [get]
func (s *Server) GetSubscriptions(c *gin.Context) {
	subscriptions, err := s.subscriptionUsecase.List()
	if err != nil {
		s.respondWithGRPCError(c, err, "failed to list subscriptions")
		return
	}

	response := make([]dto.Subscription, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		response = append(response, dto.Subscription{
			ID:        subscription.ID,
			Owner:     subscription.Owner,
			RepoName:  subscription.RepoName,
			CreatedAt: subscription.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.SubscriptionsResponse{
		Subscriptions: response,
	})
}

// GetSubscriptionsRepoInfo godoc
// @Summary Get subscriptions repositories info
// @Description Get subscriptions repositories info
// @Tags subscriptions
// @Produce json
// @Success 200 {object} dto.SubscriptionsRepoInfoResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 503 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/info [get]
func (s *Server) GetSubscriptionsRepoInfo(c *gin.Context) {
	repositories, err := s.subscriptionsRepoInfoUsecase.Execute()
	if err != nil {
		s.respondWithGRPCError(c, err, "failed to get subscriptions info")
		return
	}

	response := make([]dto.SubscriptionRepoInfo, 0, len(repositories))
	for _, repository := range repositories {
		response = append(response, dto.SubscriptionRepoInfo{
			Name:    repository.Name,
			Description: repository.Description,
			Stars:       int64(repository.StarsCount),
			Forks:       int64(repository.ForksCount),
			CreatedAt:   repository.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.SubscriptionsRepoInfoResponse{
		Repositories: response,
	})
}

func (s *Server) respondWithGRPCError(c *gin.Context, err error, message string) {
	s.log.Error(message, "error", err)
	c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: message})
}
