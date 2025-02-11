-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short VARCHAR NOT NULL UNIQUE,
    long VARCHAR NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE urls;
-- +goose StatementEnd
