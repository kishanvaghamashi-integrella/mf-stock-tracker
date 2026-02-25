-- +goose Up
CREATE TABLE assets (
    id BIGSERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    instrument_type TEXT NOT NULL CHECK (
        instrument_type IN ('stock', 'mutual_fund')
    ),
    isin TEXT UNIQUE,
    exchange TEXT,
    currency TEXT DEFAULT 'INR',
    external_platform_id TEXT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- +goose Down
DROP TABLE IF EXISTS assets;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd