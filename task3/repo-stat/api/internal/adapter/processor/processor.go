package processor

import (
	"context"
	"log/slog"
	"repo-stat/api/internal/domain"
	processorpb "repo-stat/proto/processor"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorpb.ProcessorClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   processorpb.NewProcessorClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	resp, err := c.pb.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		c.log.Error("processor ping failed", "error", err)
		return domain.PingStatusDown
	}

	if resp.GetStatus() != string(domain.PingStatusUp) {
		c.log.Error("processor ping returned non-up status", "status", resp.GetStatus())
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.pb.GetRepoInfo(ctx, &processorpb.GetRepoInfoRequest{
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

func (c *Client) Close() error {
	return c.conn.Close()
}
