-- name: CreateChirp :one
INSERT INTO chirp (
    id,
    created_at,
    updated_at, 
    body, 
    user_id
) VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;