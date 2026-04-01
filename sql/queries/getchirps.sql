-- name: GetChirps :many
SELECT * FROM chirp ORDER BY created_at ASC;