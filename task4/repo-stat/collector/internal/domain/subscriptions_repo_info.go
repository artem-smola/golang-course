package domain

import "time"

type Repository struct {
	Owner    string
	RepoName string
}

type SubscriptionRepoInfo struct {
	Owner       string
	RepoName    string
	Name        string
	Description string
	StarsCount  int
	ForksCount  int
	CreatedAt   time.Time
}
