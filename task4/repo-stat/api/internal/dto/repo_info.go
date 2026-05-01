package dto

type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
