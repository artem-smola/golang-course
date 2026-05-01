package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"repo-stat/collector/internal/domain"
	gen "repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Usecase interface {
	Execute(owner, repoName string) (*domain.RepoInfo, error)
}

type GRPCServer struct {
	gen.UnimplementedCollectorServer
	log     *slog.Logger
	useCase Usecase
}

func NewGRPCServer(log *slog.Logger, useCase Usecase) *GRPCServer {
	return &GRPCServer{log: log, useCase: useCase}
}

func (s *GRPCServer) GetRepoInfo(_ context.Context, req *gen.GetRepoInfoRequest) (*gen.GetRepoInfoResponse, error) {
	s.log.Debug("processor repo-info request recieved")

	repoInfo, err := s.useCase.Execute(req.GetOwner(), req.GetRepoName())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get repo info: %s", err.Error()))
	}

	return &gen.GetRepoInfoResponse{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		StarsCount:  int64(repoInfo.StarsCount),
		ForksCount:  int64(repoInfo.ForksCount),
		CreatedAt:   repoInfo.CreatedAt.Format(time.RFC3339),
	}, nil
}
