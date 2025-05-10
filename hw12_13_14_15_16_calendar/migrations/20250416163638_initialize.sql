-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       TEXT      NOT NULL,
    description TEXT,
    start_time  TIMESTAMP NOT NULL,
    end_time    TIMESTAMP NOT NULL,
    notify_time TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
