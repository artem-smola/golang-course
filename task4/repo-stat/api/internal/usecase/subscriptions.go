package usecase

import (
	"errors"

	"github.com/artem-smola/golang-course/task4/api/internal/domain"
)

var (
	ErrEmptyOwner    = errors.New("owner is required")
	ErrEmptyRepoName = errors.New("repo_name is required")
)

type SubscriptionProvider interface {
	AddSubscription(owner, repoName string) (*domain.Subscription, error)
	DeleteSubscription(owner, repoName string) error
	GetSubscriptions() ([]domain.Subscription, error)
}

type Subscription struct {
	provider SubscriptionProvider
}

func NewSubscription(provider SubscriptionProvider) *Subscription {
	return &Subscription{
		provider: provider,
	}
}

func (u *Subscription) Add(owner, repoName string) (*domain.Subscription, error) {
	if owner == "" {
		return nil, ErrEmptyOwner
	}

	if repoName == "" {
		return nil, ErrEmptyRepoName
	}

	return u.provider.AddSubscription(owner, repoName)
}

func (u *Subscription) Delete(owner, repoName string) error {
	if owner == "" {
		return ErrEmptyOwner
	}

	if repoName == "" {
		return ErrEmptyRepoName
	}

	return u.provider.DeleteSubscription(owner, repoName)
}

func (u *Subscription) List() ([]domain.Subscription, error) {
	return u.provider.GetSubscriptions()
}
