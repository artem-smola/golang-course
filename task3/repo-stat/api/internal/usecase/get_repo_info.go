package usecase

import (
	"errors"
	"net/url"
	"repo-stat/api/internal/domain"
	"strings"
)

var (
	ErrEmptyRepositoryURL   = errors.New("url query parameter is required")
	ErrInvalidRepositoryURL = errors.New("invalid github repository url")
)

type RepoInfoProvider interface {
	GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error)
}

type GetRepoInfo struct {
	provider RepoInfoProvider
}

func NewGetRepoInfo(provider RepoInfoProvider) *GetRepoInfo {
	return &GetRepoInfo{
		provider: provider,
	}
}

func (u *GetRepoInfo) Execute(rawURL string) (*domain.RepoInfo, error) {
	owner, repoName, err := parseGitHubRepoURL(rawURL)
	if err != nil {
		return nil, err
	}

	return u.provider.GetRepoInfo(owner, repoName)
}

func parseGitHubRepoURL(rawURL string) (string, string, error) {
	if strings.TrimSpace(rawURL) == "" {
		return "", "", ErrEmptyRepositoryURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", ErrInvalidRepositoryURL
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", "", ErrInvalidRepositoryURL
	}

	host := strings.ToLower(parsedURL.Host)
	if host != "github.com" && host != "www.github.com" {
		return "", "", ErrInvalidRepositoryURL
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 2 {
		return "", "", ErrInvalidRepositoryURL
	}

	owner := strings.TrimSpace(pathParts[0])
	repoName := strings.TrimSuffix(strings.TrimSpace(pathParts[1]), ".git")
	if owner == "" || repoName == "" {
		return "", "", ErrInvalidRepositoryURL
	}

	return owner, repoName, nil
}
