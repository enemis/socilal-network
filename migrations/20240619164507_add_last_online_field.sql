-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE users ADD column last_online timestamp default now();
CREATE INDEX user_last_online_indx on users (last_online)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX user_last_online_indx;
ALTER TABLE users DROP column last_online;
-- +goose StatementEnd
