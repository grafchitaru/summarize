-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users"
(
    id uuid PRIMARY KEY NOT NULL,
    created_at timestamp(0) without time zone NOT NULL,
    login text NOT NULL,
    password text NOT NULL,
    salt text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
