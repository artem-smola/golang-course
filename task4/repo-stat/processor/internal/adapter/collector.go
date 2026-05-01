package adapter

import (
	"context"
	"log/slog"
	"time"

	"github.com/artem-smola/golang-course/task4/processor/internal/domain"
	collectorpb "github.com/artem-smola/golang-course/task4/proto/collector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.CollectorClient
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
		pb:   collectorpb.NewCollectorClient(conn),
	}, nil
}

func (c *Client) GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error) {
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

func (c *Client) GetSubscriptionsRepoInfo() ([]domain.SubscriptionRepoInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	resp, err := c.pb.GetSubscriptionsRepoInfo(ctx, &collectorpb.GetSubscriptionsRepoInfoRequest{})
	if err != nil {
		return nil, err
	}

	repositoriesPb := resp.GetRepositories()
	repositories := make([]domain.SubscriptionRepoInfo, 0, len(repositoriesPb))
	for _, repository := range repositoriesPb {
		repositories = append(repositories, domain.SubscriptionRepoInfo{
			Name:        repository.GetName(),
			Description: repository.GetDescription(),
			StarsCount:  int(repository.GetStarsCount()),
			ForksCount:  int(repository.GetForksCount()),
			CreatedAt:   repository.GetCreatedAt(),
		})
	}

	return repositories, nil	
}

func (c *Client) Close() error {
	return c.conn.Close()
}
