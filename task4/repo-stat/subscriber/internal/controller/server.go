package controller

import (
	"context"
	"errors"
	"log/slog"
	"time"

	subscriberpb "github.com/artem-smola/golang-course/task4/proto/subscriber"
	"github.com/artem-smola/golang-course/task4/subscriber/internal/domain"
	"github.com/artem-smola/golang-course/task4/subscriber/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	subscriberpb.UnimplementedSubscriberServer
	log          *slog.Logger
	ping         *usecase.Ping
	subscription *usecase.Subscription
}

func NewServer(log *slog.Logger, ping *usecase.Ping, subscription *usecase.Subscription) *Server {
	return &Server{
		log:          log,
		ping:         ping,
		subscription: subscription,
	}
}

func (s *Server) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	s.log.Debug("subscriber ping request received")

	return &subscriberpb.PingResponse{
		Status: s.ping.Execute(ctx),
	}, nil
}

func (s *Server) AddSubscription(ctx context.Context, req *subscriberpb.AddSubscriptionRequest) (*subscriberpb.AddSubscriptionResponse, error) {
	subscription, err := s.subscription.Add(ctx, req.GetOwner(), req.GetRepoName())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOwner), errors.Is(err, domain.ErrInvalidRepoName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, domain.ErrRepositoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, domain.ErrSubscriptionExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, domain.ErrRepositoryCheckFailed):
			return nil, status.Error(codes.Unavailable, err.Error())
		default:
			s.log.Error("failed to add subscription", "error", err)
			return nil, status.Error(codes.Internal, "failed to add subscription")
		}
	}

	return &subscriberpb.AddSubscriptionResponse{
		Subscription: &subscriberpb.Subscription{
			Id:        subscription.ID,
			Owner:     subscription.Owner,
			RepoName:  subscription.RepoName,
			CreatedAt: subscription.CreatedAt.UTC().Format(time.RFC3339),
		},
	}, nil
}

func (s *Server) DeleteSubscription(ctx context.Context, req *subscriberpb.DeleteSubscriptionRequest) (*subscriberpb.DeleteSubscriptionResponse, error) {
	err := s.subscription.Delete(ctx, req.GetOwner(), req.GetRepoName())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOwner), errors.Is(err, domain.ErrInvalidRepoName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, domain.ErrSubscriptionNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			s.log.Error("failed to delete subscription", "error", err)
			return nil, status.Error(codes.Internal, "failed to delete subscription")
		}
	}

	return &subscriberpb.DeleteSubscriptionResponse{}, nil
}

func (s *Server) GetSubscriptions(ctx context.Context, _ *subscriberpb.GetSubscriptionsRequest) (*subscriberpb.GetSubscriptionsResponse, error) {
	subscriptions, err := s.subscription.GetSubscriptions(ctx)
	if err != nil {
		s.log.Error("failed to list subscriptions", "error", err)
		return nil, status.Error(codes.Internal, "failed to list subscriptions")
	}

	items := make([]*subscriberpb.Subscription, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		items = append(items, &subscriberpb.Subscription{
			Id:        subscription.ID,
			Owner:     subscription.Owner,
			RepoName:  subscription.RepoName,
			CreatedAt: subscription.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return &subscriberpb.GetSubscriptionsResponse{
		Subscriptions: items,
	}, nil
}

func (s *Server) GetSubscriptionsRepoInfo(ctx context.Context, _ *subscriberpb.GetSubscriptionsRepoInfoRequest) (*subscriberpb.GetSubscriptionsRepoInfoResponse, error) {
	repositories, err := s.subscription.GetSubscriptionsRepoInfo(ctx)
	if err != nil {
		s.log.Error("failed to list tracked repositories", "error", err)
		return nil, status.Error(codes.Internal, "failed to list tracked repositories")
	}

	items := make([]*subscriberpb.SubscriptionRepoInfo, 0, len(repositories))
	for _, repository := range repositories {
		items = append(items, &subscriberpb.SubscriptionRepoInfo{
			Owner:    repository.Owner,
			RepoName: repository.RepoName,
		})
	}

	return &subscriberpb.GetSubscriptionsRepoInfoResponse{
		Repositories: items,
	}, nil
}
