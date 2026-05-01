package dto

type AddSubscriptionRequest struct {
	Owner    string `json:"owner" binding:"required" example:"Desbordante"`
	RepoName string `json:"repo_name" binding:"required" example:"desbordante-core"`
}

type Subscription struct {
	ID        int64  `json:"id"`
	Owner     string `json:"owner"`
	RepoName  string `json:"repo_name"`
	CreatedAt string `json:"created_at"`
}

type SubscriptionRepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type SubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type SubscriptionsRepoInfoResponse struct {
	Repositories []SubscriptionRepoInfo `json:"repositories"`
}
