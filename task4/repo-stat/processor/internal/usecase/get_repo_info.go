package usecase

import "github.com/artem-smola/golang-course/task4/processor/internal/domain"

type Client interface {
	GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error)
	GetSubscriptionsRepoInfo() ([]domain.SubscriptionRepoInfo, error)
}

type GetRepoInfoUsecase struct {
	client Client
}

func NewGetRepoInfoUsecase(client Client) *GetRepoInfoUsecase {
	return &GetRepoInfoUsecase{client: client}
}

func (u *GetRepoInfoUsecase) Execute(owner, repoName string) (*domain.RepoInfo, error) {
	return u.client.GetRepoInfo(owner, repoName)
}
