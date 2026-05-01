-- name: CreateSubscription :one
INSERT INTO subscriptions (owner, repo_name)
VALUES ($1, $2)
RETURNING id, owner, repo_name, created_at;

-- name: DeleteSubscription :execrows
DELETE FROM subscriptions
WHERE owner = $1
    AND repo_name = $2;

-- name: GetSubscriptions :many
SELECT id, owner, repo_name, created_at
FROM subscriptions
ORDER BY id;

-- name: GetSubscriptionsRepoInfo :many
SELECT owner, repo_name
FROM subscriptions
ORDER BY id;