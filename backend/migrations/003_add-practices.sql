-- +goose Up
CREATE TABLE practices (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    note_id      UUID NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    status       TEXT NOT NULL DEFAULT 'in_progress'
                   CHECK (status IN ('pending', 'in_progress', 'completed')),
    exercises    JSONB NOT NULL DEFAULT '[]',
    created_at   TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'),
    completed_at TIMESTAMP
);

CREATE INDEX practices_note_id_idx    ON practices(note_id);
CREATE INDEX practices_note_status_idx ON practices(note_id, status);

-- +goose Down
DROP TABLE practices;
