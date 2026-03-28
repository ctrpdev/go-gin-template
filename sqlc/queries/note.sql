-- name: CreateNote :one
INSERT INTO notes (user_id, title, content)
VALUES ($1, $2, $3)
RETURNING id, user_id, title, content, created_at, updated_at;

-- name: DeleteNote :exec
DELETE FROM notes
WHERE id = $1;

-- name: GetNoteByID :one
SELECT *
FROM notes
WHERE id = $1 LIMIT 1;

-- name: ListNotes :many
SELECT *
FROM notes;

-- name: ListUserNotes :many
SELECT *
FROM notes
WHERE user_id = $1;

-- name: UpdateNote :one
UPDATE notes
SET title = $2, content = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, title, content, created_at, updated_at;