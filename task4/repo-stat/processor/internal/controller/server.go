package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/artem-smola/golang-course/task4/processor/internal/domain"
	"github.com/artem-smola/golang-course/task4/processor/internal/usecase"
	processorpb "github.com/artem-smola/golang-course/task4/proto/processor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetRepoInfoUsecase interface {
	Execute(owner, repoName string) (*domain.RepoInfo, error)
}

type Server struct {
	processorpb.UnimplementedProcessorServer
	log                  *slog.Logger
	getRepoInfo          *usecase.GetRepoInfoUsecase
	getSubscriptionsRepoInfo *usecase.GetSubscriptionsRepoInfoUsecase
	ping                 *usecase.Ping
}

func NewServer(log *slog.Logger, getRepoInfo *usecase.GetRepoInfoUsecase, getSubscriptionsRepoInfo *usecase.GetSubscriptionsRepoInfoUsecase, ping *usecase.Ping) *Server {
	return &Server{
		log:                  log,
		getRepoInfo:          getRepoInfo,
		getSubscriptionsRepoInfo: getSubscriptionsRepoInfo,
		ping:                 ping,
	}
}

func (s *Server) Ping(_ context.Context, _ *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	s.log.Debug("processorpb ping request received")

	return &processorpb.PingResponse{
		Status: s.ping.Execute(),
	}, nil
}

func (s *Server) GetRepoInfo(_ context.Context, req *processorpb.GetRepoInfoRequest) (*processorpb.GetRepoInfoResponse, error) {
	s.log.Debug("processorpb get-repo-info request recieved")

	repoInfo, err := s.getRepoInfo.Execute(req.GetOwner(), req.GetRepoName())
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get repo info: %s", err.Error()))
	}

	return &processorpb.GetRepoInfoResponse{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		StarsCount:  int64(repoInfo.StarsCount),
		ForksCount:  int64(repoInfo.ForksCount),
		CreatedAt:   repoInfo.CreatedAt,
	}, nil
}

func (s *Server) GetSubscriptionsRepoInfo(_ context.Context, _ *processorpb.GetSubscriptionsRepoInfoRequest) (*processorpb.GetSubscriptionsRepoInfoResponse, error) {
	s.log.Debug("processorpb get-subscriptions-info request recieved")

	repositories, err := s.getSubscriptionsRepoInfo.Execute()
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, status.Error(codes.DeadlineExceeded, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get subscriptions info: %s", err.Error()))
	}

	items := make([]*processorpb.SubscriptionRepoInfo, 0, len(repositories))
	for _, repository := range repositories {
		items = append(items, &processorpb.SubscriptionRepoInfo{
			Name:        repository.Name,
			Description: repository.Description,
			StarsCount:  int64(repository.StarsCount),
			ForksCount:  int64(repository.ForksCount),
			CreatedAt:   repository.CreatedAt,
		})
	}

	return &processorpb.GetSubscriptionsRepoInfoResponse{
		Repositories: items,
	}, nil
}
