package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/artem-smola/golang-course/task4/collector/internal/domain"
	collectorpb "github.com/artem-smola/golang-course/task4/proto/collector"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecase interface {
	GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error)
	GetSubscriptionsRepoInfo(ctx context.Context) ([]domain.SubscriptionRepoInfo, error)
}

type Server struct {
	collectorpb.UnimplementedCollectorServer
	log     *slog.Logger
	useCase Usecase
}

func NewServer(log *slog.Logger, useCase Usecase) *Server {
	return &Server{log: log, useCase: useCase}
}

func (s *Server) GetRepoInfo(_ context.Context, req *collectorpb.GetRepoInfoRequest) (*collectorpb.GetRepoInfoResponse, error) {
	s.log.Debug("processor repo-info request recieved")

	repoInfo, err := s.useCase.GetRepoInfo(req.GetOwner(), req.GetRepoName())
	if err != nil {
		if errors.Is(err, domain.ErrRepositoryNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get repo info: %s", err.Error()))
	}

	return &collectorpb.GetRepoInfoResponse{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		StarsCount:  int64(repoInfo.StarsCount),
		ForksCount:  int64(repoInfo.ForksCount),
		CreatedAt:   repoInfo.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *Server) GetSubscriptionsRepoInfo(ctx context.Context, _ *collectorpb.GetSubscriptionsRepoInfoRequest) (*collectorpb.GetSubscriptionsRepoInfoResponse, error) {
	s.log.Debug("processor subscriptions-info request recieved")

	repositories, err := s.useCase.GetSubscriptionsRepoInfo(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrRepositoryNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get subscriptions info: %s", err.Error()))
	}

	items := make([]*collectorpb.SubscriptionRepoInfo, 0, len(repositories))
	for _, repository := range repositories {
		items = append(items, &collectorpb.SubscriptionRepoInfo{
			Name:        repository.Name,
			Description: repository.Description,
			StarsCount:  int64(repository.StarsCount),
			ForksCount:  int64(repository.ForksCount),
			CreatedAt:   repository.CreatedAt.Format(time.RFC3339),
		})
	}

	return &collectorpb.GetSubscriptionsRepoInfoResponse{
		Repositories: items,
	}, nil
}