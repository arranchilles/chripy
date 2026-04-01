-- name: GetUserFromRefreshToken :one
SELECT * FROM users JOIN refresh_tokens ON user_id = id WHERE token = $1;