package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/artem-smola/golang-course/task4/subscriber/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, owner, repoName string) (*domain.Subscription, error)
	Delete(ctx context.Context, owner, repoName string) error
	GetSubscriptions(ctx context.Context) ([]domain.Subscription, error)
	GetSubscriptionsRepoInfo(ctx context.Context) ([]domain.SubscriptionRepoInfo, error)
}

type RepositoryValidator interface {
	ValidateExists(ctx context.Context, owner, repoName string) error
}

type Subscription struct {
	repository SubscriptionRepository
	validator  RepositoryValidator
}

func NewSubscription(repository SubscriptionRepository, validator RepositoryValidator) *Subscription {
	return &Subscription{
		repository: repository,
		validator:  validator,
	}
}

func (u *Subscription) Add(ctx context.Context, owner, repoName string) (*domain.Subscription, error) {
	owner = strings.TrimSpace(owner)
	if owner == "" {
		return nil, domain.ErrInvalidOwner
	}

	repoName = normalizeRepoName(repoName)
	if repoName == "" {
		return nil, domain.ErrInvalidRepoName
	}

	if err := u.validator.ValidateExists(ctx, owner, repoName); err != nil {
		if errors.Is(err, domain.ErrRepositoryNotFound) {
			return nil, err
		}
		return nil, domain.ErrRepositoryCheckFailed
	}

	subscription, err := u.repository.Create(ctx, owner, repoName)
	if err != nil {
		if errors.Is(err, domain.ErrSubscriptionExists) {
			return nil, err
		}
		return nil, domain.ErrSubscriptionStoreFailed
	}

	return subscription, nil
}

func (u *Subscription) Delete(ctx context.Context, owner, repoName string) error {
	owner = strings.TrimSpace(owner)
	if owner == "" {
		return domain.ErrInvalidOwner
	}

	repoName = normalizeRepoName(repoName)
	if repoName == "" {
		return domain.ErrInvalidRepoName
	}

	if err := u.repository.Delete(ctx, owner, repoName); err != nil {
		if errors.Is(err, domain.ErrSubscriptionNotFound) {
			return err
		}
		return domain.ErrSubscriptionStoreFailed
	}

	return nil
}

func (u *Subscription) GetSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	subscriptions, err := u.repository.GetSubscriptions(ctx)
	if err != nil {
		return nil, domain.ErrSubscriptionStoreFailed
	}
	return subscriptions, nil
}

func (u *Subscription) GetSubscriptionsRepoInfo(ctx context.Context) ([]domain.SubscriptionRepoInfo, error) {
	repositories, err := u.repository.GetSubscriptionsRepoInfo(ctx)
	if err != nil {
		return nil, domain.ErrSubscriptionStoreFailed
	}
	return repositories, nil
}


func normalizeRepoName(repoName string) string {
	repoName = strings.TrimSpace(repoName)
	return strings.TrimSuffix(repoName, ".git")
}
