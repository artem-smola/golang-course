package usecase

import "repo-stat/processor/internal/domain"

type GRPCClient interface {
	GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error)
}

type GetRepoInfo struct {
	client GRPCClient
}

func NewGetRepoInfo(client GRPCClient) *GetRepoInfo {
	return &GetRepoInfo{client: client}
}

func (u *GetRepoInfo) Execute(owner, repoName string) (*domain.RepoInfo, error) {
	return u.client.GetRepoInfo(owner, repoName)
}
