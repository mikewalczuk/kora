-- name: CreatePractice :one
INSERT INTO practices (note_id, status, exercises)
VALUES (@note_id, @status, @exercises)
RETURNING *;

-- name: GetPractice :one
SELECT * FROM practices WHERE id = @id;

-- name: GetActivePracticeForNote :one
SELECT * FROM practices
WHERE note_id = @note_id
  AND status IN ('pending', 'in_progress')
LIMIT 1;

-- name: ListPractices :many
SELECT * FROM practices
WHERE (sqlc.narg('note_id')::uuid IS NULL OR note_id = sqlc.narg('note_id'))
  AND (cardinality(sqlc.arg('statuses')::text[]) = 0 OR status = ANY(sqlc.arg('statuses')::text[]))
ORDER BY created_at DESC
LIMIT sqlc.arg('lim')
OFFSET sqlc.arg('off');

-- name: CountPractices :one
SELECT COUNT(*) FROM practices
WHERE (sqlc.narg('note_id')::uuid IS NULL OR note_id = sqlc.narg('note_id'))
  AND (cardinality(sqlc.arg('statuses')::text[]) = 0 OR status = ANY(sqlc.arg('statuses')::text[]));

-- name: UpdatePracticeExercises :one
UPDATE practices SET exercises = @exercises WHERE id = @id RETURNING *;

-- name: CompletePractice :one
UPDATE practices
SET status = 'completed', completed_at = (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')
WHERE id = @id
RETURNING *;
