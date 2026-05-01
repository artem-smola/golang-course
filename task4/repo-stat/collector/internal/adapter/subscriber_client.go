package adapter

import (
	"context"
	"log/slog"
	"time"

	"github.com/artem-smola/golang-course/task4/collector/internal/domain"
	subscriberpb "github.com/artem-smola/golang-course/task4/proto/subscriber"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SubscriberClient struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscriberpb.SubscriberClient
}

func NewSubscriberClient(address string, log *slog.Logger) (*SubscriberClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &SubscriberClient{
		log:  log,
		conn: conn,
		pb:   subscriberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *SubscriberClient) GetSubscriptionsRepoInfo(ctx context.Context) ([]domain.Repository, error) {
	requestCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	resp, err := c.pb.GetSubscriptionsRepoInfo(requestCtx, &subscriberpb.GetSubscriptionsRepoInfoRequest{})
	if err != nil {
		return nil, err
	}

	repositories := make([]domain.Repository, 0, len(resp.GetRepositories()))
	for _, repository := range resp.GetRepositories() {
		repositories = append(repositories, domain.Repository{
			Owner:    repository.GetOwner(),
			RepoName: repository.GetRepoName(),
		})
	}

	return repositories, nil
}

func (c *SubscriberClient) Close() error {
	return c.conn.Close()
}
