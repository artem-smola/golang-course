package adapter

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/artem-smola/golang-course/task4/subscriber/internal/domain"
)

type Client struct {
	url    string
	token  string
	client *http.Client
}

func NewClient(url, token string, timeout time.Duration) *Client {
	url = strings.TrimRight(strings.TrimSpace(url), "/")
	if url == "" {
		url = "https://api.github.com"
	}

	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		url:   url,
		token: strings.TrimSpace(token),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) ValidateExists(ctx context.Context, owner, repoName string) error {
	repoURL := fmt.Sprintf("%s/repos/%s/%s", c.url, url.PathEscape(owner), url.PathEscape(repoName))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, repoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "repo-stat-subscriber")
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return domain.ErrRepositoryNotFound
	default:
		return fmt.Errorf("github api status: %d", resp.StatusCode)
	}
}
