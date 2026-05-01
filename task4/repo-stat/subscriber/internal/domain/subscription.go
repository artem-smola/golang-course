package domain

import "time"

type Subscription struct {
	ID        int64
	Owner     string
	RepoName  string
	CreatedAt time.Time
}

type SubscriptionRepoInfo struct {
	Owner    string
	RepoName string
}
