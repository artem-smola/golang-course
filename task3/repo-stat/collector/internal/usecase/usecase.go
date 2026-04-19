package usecase

import "repo-stat/collector/internal/domain"

type HTTPClient interface {
	GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error)
}

type Usecase struct {
	client HTTPClient
}

func NewUsecase(client HTTPClient) *Usecase {
	return &Usecase{client: client}
}

func (u *Usecase) Execute(owner, repoName string) (*domain.RepoInfo, error) {
	return u.client.GetRepoInfo(owner, repoName)
}
