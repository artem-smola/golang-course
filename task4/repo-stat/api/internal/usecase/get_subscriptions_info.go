package usecase

import "github.com/artem-smola/golang-course/task4/api/internal/domain"

type SubscriptionsRepoInfoProvider interface {
	GetSubscriptionsRepoInfo() ([]domain.SubscriptionRepoInfo, error)
}

type GetSubscriptionsRepoInfoUsecase struct {
	provider SubscriptionsRepoInfoProvider
}

func NewGetSubscriptionsRepoInfoUsecase(provider SubscriptionsRepoInfoProvider) *GetSubscriptionsRepoInfoUsecase {
	return &GetSubscriptionsRepoInfoUsecase{
		provider: provider,
	}
}

func (u *GetSubscriptionsRepoInfoUsecase) Execute() ([]domain.SubscriptionRepoInfo, error) {
	return u.provider.GetSubscriptionsRepoInfo()
}
