package adapter

import (
	"context"
	"log/slog"
	"time"

	"repo-stat/processor/internal/domain"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Создаем grpc клиента, посылающего запрос коллектору

type GRPCClient struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.CollectorClient
}

func NewGRPCClient(address string, log *slog.Logger) (*GRPCClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		log:  log,
		conn: conn,
		pb:   collectorpb.NewCollectorClient(conn),
	}, nil
}

func (c *GRPCClient) GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	resp, err := c.pb.GetRepoInfo(ctx, &collectorpb.GetRepoInfoRequest{
		Owner:    owner,
		RepoName: repoName,
	})
	if err != nil {
		return nil, err
	}

	return &domain.RepoInfo{
		Name:        resp.GetName(),
		Description: resp.GetDescription(),
		StarsCount:  int(resp.GetStarsCount()),
		ForksCount:  int(resp.GetForksCount()),
		CreatedAt:   resp.GetCreatedAt(),
	}, nil
}
