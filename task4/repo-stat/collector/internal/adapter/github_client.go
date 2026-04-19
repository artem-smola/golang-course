package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/artem-smola/golang-course/task4/collector/internal/domain"
)

type RepoInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StarsCount  int       `json:"stargazers_count"`
	ForksCount  int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type GitHubClient struct{}

func (c *GitHubClient) GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repoName)

	const requestTimeout = 10 * time.Second
	const maxAttempts = 5
	client := &http.Client{Timeout: requestTimeout}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		repoInfo, err := c.fetchRepoInfo(client, url)
		if err == nil {
			name := repoInfo.Name
			if name == "" {
				name = fmt.Sprintf("%s/%s", owner, repoName)
			}

			return &domain.RepoInfo{
				Name:        name,
				Description: repoInfo.Description,
				StarsCount:  repoInfo.StarsCount,
				ForksCount:  repoInfo.ForksCount,
				CreatedAt:   repoInfo.CreatedAt,
			}, nil
		}

		lastErr = err
		if !isRetriable(err) || attempt == maxAttempts {
			break
		}

		time.Sleep(time.Duration(attempt) * 300 * time.Millisecond)
	}

	return nil, lastErr
}

func (c *GitHubClient) fetchRepoInfo(client *http.Client, url string) (*RepoInfo, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "repo-stat-collector")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, domain.ErrRepositoryNotFound
		}
		return nil, fmt.Errorf("github api status: %d", resp.StatusCode)
	}

	var repoInfo RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		return nil, err
	}
	return &repoInfo, nil
}

func isRetriable(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	msg := err.Error()
	if errors.Is(err, domain.ErrRepositoryNotFound) {
		return false
	}
	if containsAny(msg, "timeout", "tls", "connection refused", "temporary", "github api status: 5", "github api status: 429") {
		return true
	}
	return false
}

func containsAny(s string, needles ...string) bool {
	for _, needle := range needles {
		if needle != "" && strings.Contains(s, needle) {
			return true
		}
	}
	return false
}
