-- name: InsertEntry :one
INSERT INTO entries(
    account_id,
    amount
) VALUES(
    $1, $2
) RETURNING *;

-- name: GetEntryById :one
SELECT * FROM entries
WHERE id = $1
LIMIT 1;

-- name: GetEntriesForAccount :many
SELECT * FROM entries
WHERE account_id = $1 
LIMIT $2
OFFSET $3;

-- name: UpdateEntry :one
UPDATE entries
SET amount = $1
WHERE id = $2
RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;