-- name: UpdateUserToken :one
UPDATE users
SET authtoken = $1
WHERE name = $2
RETURNING *;