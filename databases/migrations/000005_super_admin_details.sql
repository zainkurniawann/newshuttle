-- +goose Up
-- +goose StatementBegin
CREATE TABLE super_admin_details (
    user_uuid UUID PRIMARY KEY,
    user_picture TEXT NULL,
    user_first_name VARCHAR(100),
    user_last_name VARCHAR(100),
    user_gender VARCHAR(20),
    user_phone VARCHAR(50),
    user_address TEXT,
    FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS super_admin_details;
-- +goose StatementEnd
