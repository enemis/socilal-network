-- +goose NO TRANSACTION
-- +goose Up
CREATE INDEX CONCURRENTLY IF NOT EXISTS name_index ON users (lower((name)::text) varchar_pattern_ops);
CREATE INDEX CONCURRENTLY IF NOT EXISTS surname_index ON users (lower((surname)::text) varchar_pattern_ops);

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS name_index;
DROP INDEX IF EXISTS surname_index;
-- +goose StatementEnd
