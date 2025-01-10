-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS fcm_tokens (
    id BIGINT PRIMARY KEY,
    user_uuid UUID UNIQUE NOT NULL,
    device_token TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NULL DEFAULT NULL,
    FOREIGN KEY (user_uuid) REFERENCES users (user_uuid) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fcm_tokens CASCADE;
-- +goose StatementEnd
