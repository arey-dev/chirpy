-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * from users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET
  email = $1,
  hashed_password = $2,
  updated_at = now()
WHERE id = $3
RETURNING *;

-- name: UpdateUserToChirpyRed :one
UPDATE users
SET 
  is_chirpy_red = true
where id = $1
RETURNING *;
