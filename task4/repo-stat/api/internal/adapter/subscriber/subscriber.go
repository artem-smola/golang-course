package subscriber

import (
	"context"
	"log/slog"
	"time"

	subscriberpb "github.com/artem-smola/golang-course/task4/proto/subscriber"

	"github.com/artem-smola/golang-course/task4/api/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscriberpb.SubscriberClient
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
		pb:   subscriberpb.NewSubscriberClient(conn),
	}, nil
}

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	resp, err := c.pb.Ping(ctx, &subscriberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", "error", err)
		return domain.PingStatusDown
	}
	if resp.GetStatus() != string(domain.PingStatusUp) {
		c.log.Error("subscriber ping returned non-up status", "status", resp.GetStatus())
		return domain.PingStatusDown
	}

	return domain.PingStatusUp
}

func (c *Client) AddSubscription(owner, repoName string) (*domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.pb.AddSubscription(ctx, &subscriberpb.AddSubscriptionRequest{
		Owner: owner,
		RepoName: repoName,
	})
	if err != nil {
		return nil, err
	}

	subscription := resp.GetSubscription()
	return &domain.Subscription{
		ID: subscription.GetId(),
		Owner: subscription.GetOwner(),
		RepoName: subscription.GetRepoName(),
		CreatedAt: subscription.GetCreatedAt(),
	}, nil
}

func (c *Client) DeleteSubscription(owner, repoName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := c.pb.DeleteSubscription(ctx, &subscriberpb.DeleteSubscriptionRequest{
		Owner:    owner,
		RepoName: repoName,
	})
	return err
}

func (c *Client) GetSubscriptions() ([]domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.pb.GetSubscriptions(ctx, &subscriberpb.GetSubscriptionsRequest{})
	if err != nil {
		return nil, err
	}

	subscriptions := resp.GetSubscriptions()
	result := make([]domain.Subscription, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		result = append(result, domain.Subscription{
			ID: subscription.GetId(),
			Owner: subscription.GetOwner(),
			RepoName: subscription.GetRepoName(),
			CreatedAt: subscription.GetCreatedAt(),
		})
	}
	return result, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

