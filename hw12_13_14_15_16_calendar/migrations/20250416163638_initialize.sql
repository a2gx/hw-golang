-- +goose Up
-- +goose StatementBegin
CREATE TABLE events
(
    id          UUID PRIMARY KEY,
    title       TEXT      NOT NULL,
    start_time  TIMESTAMP NOT NULL,
    end_time    TIMESTAMP NOT NULL,
    description TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
