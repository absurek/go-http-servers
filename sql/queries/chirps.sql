-- name: CreateChirp :one
INSERT INTO chirps (id, user_id, body, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, DEFAULT, DEFAULT)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at;

-- name: GetChirpByID :one
SELECT * FROM chirps WHERE id = $1;

-- name: CheckChirpAccess :one
SELECT EXISTS (
    SELECT 1
    FROM chirps
    WHERE id = $1 AND user_id = $2
);

-- name: DeleteChirp :execrows
DELETE FROM chirps WHERE id = $1 AND user_id = $2;
