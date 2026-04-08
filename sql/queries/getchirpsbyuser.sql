-- name: GetChirpsByAuthor :many
SELECT * FROM chirp WHERE user_id = $1 ORDER BY created_at ASC;