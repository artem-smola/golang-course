package domain

import "errors"

var (
	ErrInvalidOwner            = errors.New("owner is required")
	ErrInvalidRepoName         = errors.New("repo_name is required")
	ErrSubscriptionExists      = errors.New("subscription already exists")
	ErrSubscriptionNotFound    = errors.New("subscription not found")
	ErrRepositoryNotFound      = errors.New("repository not found on github")
	ErrRepositoryCheckFailed   = errors.New("failed to validate repository on github")
	ErrSubscriptionStoreFailed = errors.New("failed to persist subscription")
)
