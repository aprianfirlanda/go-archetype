-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks
(
    id          BIGSERIAL PRIMARY KEY,

    public_id   UUID      NOT NULL,
    title       TEXT      NOT NULL,
    description TEXT      NOT NULL DEFAULT '',
    status      TEXT      NOT NULL,
    priority    INT       NOT NULL,
    due_date    TIMESTAMP NULL,
    tags        TEXT      NOT NULL DEFAULT '',
    completed   BOOLEAN   NOT NULL DEFAULT FALSE,

    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_tasks_public_id UNIQUE (public_id),
    CONSTRAINT chk_tasks_priority CHECK (priority BETWEEN 1 AND 5)
);

-- Indexes for common queries
CREATE INDEX idx_tasks_status ON tasks (status);
CREATE INDEX idx_tasks_due_date ON tasks (due_date);
CREATE INDEX idx_tasks_completed ON tasks (completed);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
