-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY,
    user_uuid UUID UNIQUE NOT NULL,
    user_username VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    user_password VARCHAR(255) NOT NULL,
    user_role VARCHAR(20) NOT NULL,
    user_role_code VARCHAR (5),
    user_status VARCHAR(20) DEFAULT 'offline',
    user_last_active TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255)
);

CREATE INDEX idx_user_uuid ON users(user_uuid);
CREATE INDEX idx_user_username ON users(user_username);
CREATE INDEX idx_user_email ON users(user_email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd