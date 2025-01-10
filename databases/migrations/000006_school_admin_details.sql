-- +goose Up
-- +goose StatementBegin
CREATE TABLE school_admin_details (
    user_uuid UUID PRIMARY KEY,
    school_uuid UUID NOT NULL,
    user_picture TEXT,
    user_first_name VARCHAR(100),
    user_last_name VARCHAR(100),
    user_gender VARCHAR(20),
    user_phone VARCHAR(50),
    user_address TEXT,
    FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE,
    FOREIGN KEY (school_uuid) REFERENCES schools(school_uuid) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS school_admin_details;
-- +goose StatementEnd
