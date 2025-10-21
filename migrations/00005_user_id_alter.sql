-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS workouts
ADD COLUMN IF NOT EXISTS user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS workouts
DROP COLUMN IF EXISTS user_id;
-- +goose StatementEnd