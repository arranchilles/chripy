-- name: AddUserToChirpRed :exec
UPDATE users SET is_chirp_red = TRUE WHERE id = $1;