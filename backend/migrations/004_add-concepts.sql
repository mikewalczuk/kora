-- +goose Up
CREATE TABLE concepts (
  id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
  note_id          UUID        NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
  title            TEXT        NOT NULL,
  content          TEXT        NOT NULL,
  stability        NUMERIC     DEFAULT 0,
  difficulty       NUMERIC     DEFAULT 0,
  due_at           TIMESTAMPTZ,
  last_reviewed_at TIMESTAMPTZ,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  archived_at      TIMESTAMPTZ,
  archived_reason  TEXT        CHECK (archived_reason IN ('ai_obsoleted', 'user_archived'))
);
CREATE INDEX concepts_note_id_idx    ON concepts (note_id);
CREATE INDEX concepts_due_at_idx     ON concepts (due_at)  WHERE due_at IS NOT NULL;
CREATE INDEX concepts_active_due_idx ON concepts (due_at)  WHERE archived_at IS NULL;

-- +goose Down
DROP TABLE concepts;
