package usecase

import (
	"context"
	"fmt"

	"github.com/artem-smola/golang-course/task4/collector/internal/domain"
)

type GitHubClient interface {
	GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error)
}

type SubscriberClient interface {
	GetSubscriptionsRepoInfo(ctx context.Context) ([]domain.Repository, error)
}

type Usecase struct {
	gitHubClient     GitHubClient
	subscriberClient SubscriberClient
}

func NewUsecase(gitHubClient GitHubClient, subscriberClient SubscriberClient) *Usecase {
	return &Usecase{
		gitHubClient:     gitHubClient,
		subscriberClient: subscriberClient,
	}
}

func (u *Usecase) GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error) {
	return u.gitHubClient.GetRepoInfo(owner, repoName)
}

func (u *Usecase) GetSubscriptionsRepoInfo(ctx context.Context) ([]domain.SubscriptionRepoInfo, error) {
	repositories, err := u.subscriberClient.GetSubscriptionsRepoInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("get tracked repositories: %w", err)
	}

	result := make([]domain.SubscriptionRepoInfo, 0, len(repositories))
	for _, repository := range repositories {
		repoInfo, err := u.gitHubClient.GetRepoInfo(repository.Owner, repository.RepoName)
		if err != nil {
			return nil, fmt.Errorf("get repo info for %s/%s: %w", repository.Owner, repository.RepoName, err)
		}

		result = append(result, domain.SubscriptionRepoInfo{
			Owner:       repository.Owner,
			RepoName:    repository.RepoName,
			Name:        repoInfo.Name,
			Description: repoInfo.Description,
			StarsCount:  repoInfo.StarsCount,
			ForksCount:  repoInfo.ForksCount,
			CreatedAt:   repoInfo.CreatedAt,
		})
	}

	return result, nil
}
