-- name: CreateComment :one
INSERT INTO comments (
    comment_destination, parent_id, content, author
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetCommentsByDestination :many
SELECT * FROM comments
WHERE comment_destination = $1
ORDER BY created_at;

-- name: GetRootCommentIDsWithPagination :many
SELECT id FROM comments
WHERE comment_destination = $1
  AND parent_id IS NULL
  AND id > $2
ORDER BY id ASC
    LIMIT $3;

-- name: GetCommentsWithChildren :many
WITH RECURSIVE comment_tree AS (
    SELECT id, parent_id, comment_destination, author, content, created_at, updated_at
    FROM comments
    WHERE id = ANY(@ids::int[])

    UNION ALL

    SELECT c.id, c.parent_id, c.comment_destination, c.author, c.content, c.created_at, c.updated_at
    FROM comments c
             INNER JOIN comment_tree ct ON c.parent_id = ct.id
)
SELECT * FROM comment_tree
ORDER BY created_at ASC;

-- name: DeleteCommentByID :exec
DELETE FROM comments
WHERE id = $1;