package postgres

import (
	"context"
	"errors"

	"github.com/artem-smola/golang-course/task4/subscriber/internal/domain"
	db "github.com/artem-smola/golang-course/task4/subscriber/internal/storage/postgres/sqlc"

	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) *Repository {
	return &Repository{
		queries: queries,
	}
}

func (r *Repository) Create(ctx context.Context, owner, repoName string) (*domain.Subscription, error) {
	row, err := r.queries.CreateSubscription(ctx, db.CreateSubscriptionParams{
		Owner:    owner,
		RepoName: repoName,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrSubscriptionExists
		}
		return nil, err
	}

	createdAt := row.CreatedAt.Time
	if row.CreatedAt.Valid {
		createdAt = createdAt.UTC()
	}

	return &domain.Subscription{
		ID:        row.ID,
		Owner:     row.Owner,
		RepoName:  row.RepoName,
		CreatedAt: createdAt,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, owner, repoName string) error {
	rowsAffected, err := r.queries.DeleteSubscription(ctx, db.DeleteSubscriptionParams{
		Owner:    owner,
		RepoName: repoName,
	})
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrSubscriptionNotFound
	}
	return nil
}

func (r *Repository) GetSubscriptions(ctx context.Context) ([]domain.Subscription, error) {
	rows, err := r.queries.GetSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Subscription, 0, len(rows))
	for _, row := range rows {
		createdAt := row.CreatedAt.Time
		if row.CreatedAt.Valid {
			createdAt = createdAt.UTC()
		}

		result = append(result, domain.Subscription{
			ID:        row.ID,
			Owner:     row.Owner,
			RepoName:  row.RepoName,
			CreatedAt: createdAt,
		})
	}
	return result, nil
}

func (r *Repository) GetSubscriptionsRepoInfo(ctx context.Context) ([]domain.SubscriptionRepoInfo, error) {
	rows, err := r.queries.GetSubscriptionsRepoInfo(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.SubscriptionRepoInfo, 0, len(rows))
	for _, row := range rows {
		result = append(result, domain.SubscriptionRepoInfo{
			Owner:    row.Owner,
			RepoName: row.RepoName,
		})
	}
	return result, nil
}
