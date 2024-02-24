-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "summarize"
(
    id uuid PRIMARY KEY NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp(0) without time zone NOT NULL,
    text text NOT NULL,
    result text,
    status varchar(55) NOT NULL,
    tokens integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE summarize;
-- +goose StatementEnd
