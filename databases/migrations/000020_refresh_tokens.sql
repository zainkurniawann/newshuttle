-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    id BIGINT PRIMARY KEY,
    user_uuid UUID NOT NULL,
    refresh_token TEXT NOT NULL,
    issued_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    last_used_at TIMESTAMPTZ,
    CONSTRAINT token_fk_user FOREIGN KEY(user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE,
    CONSTRAINT unique_user_uuid UNIQUE (user_uuid)

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_tokens;
-- +goose StatementEnd