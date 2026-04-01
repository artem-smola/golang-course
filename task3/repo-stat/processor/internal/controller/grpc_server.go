package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"repo-stat/processor/internal/domain"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetRepoInfo interface {
	Execute(owner, repoName string) (*domain.RepoInfo, error)
}

type GRPCServer struct {
	processorpb.UnimplementedProcessorServer
	log         *slog.Logger
	getRepoInfo *usecase.GetRepoInfo
	ping        *usecase.Ping
}

func NewGRPCServer(log *slog.Logger, getRepoInfo *usecase.GetRepoInfo, ping *usecase.Ping) *GRPCServer {
	return &GRPCServer{log: log, getRepoInfo: getRepoInfo, ping: ping}
}

func (s *GRPCServer) Ping(_ context.Context, _ *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	s.log.Debug("processorpb ping request received")

	return &processorpb.PingResponse{
		Status: s.ping.Execute(),
	}, nil
}

func (s *GRPCServer) GetRepoInfo(_ context.Context, req *processorpb.GetRepoInfoRequest) (*processorpb.GetRepoInfoResponse, error) {
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
