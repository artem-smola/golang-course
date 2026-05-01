package domain

type Subscription struct {
	ID        int64
	Owner     string
	RepoName  string
	CreatedAt string
}

type SubscriptionRepoInfo struct {
	Name        string
	Description string
	StarsCount  int
	ForksCount  int
	CreatedAt   string
}
