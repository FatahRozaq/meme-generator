-- name: ListActivePlans :many
SELECT * FROM plans WHERE active = true ORDER BY price_idr ASC;

-- name: GetPlanByCode :one
SELECT * FROM plans WHERE code = $1 AND active = true;

-- name: GetActiveSubscriptionByUser :one
SELECT s.*, p.code AS plan_code, p.name AS plan_name, p.max_memes
FROM subscriptions s
JOIN plans p ON p.id = s.plan_id
WHERE s.user_id = $1 AND s.status = 'active' AND s.ends_at > now()
ORDER BY s.ends_at DESC LIMIT 1;

-- name: CountMemesSince :one
SELECT COUNT(*) FROM memes WHERE user_id = $1 AND created_at >= $2;
