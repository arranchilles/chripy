-- name: GetChirp :one
SELECT * FROM chirp WHERE id = $1;