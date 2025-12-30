CREATE TABLE tasks
(
    id          BIGSERIAL PRIMARY KEY,
    public_id   UUID      NOT NULL UNIQUE,

    title       TEXT      NOT NULL,
    description TEXT,
    status      TEXT      NOT NULL,
    priority    INT       NOT NULL,
    due_date    TIMESTAMP NULL,
    tags        TEXT,
    completed   BOOLEAN   NOT NULL DEFAULT FALSE,

    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);

CREATE INDEX idx_tasks_public_id ON tasks (public_id);
