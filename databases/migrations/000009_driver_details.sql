-- +goose Up
-- +goose StatementBegin
CREATE TABLE driver_details (
    user_uuid UUID PRIMARY KEY,
    school_uuid UUID NULL,
    vehicle_uuid UUID NULL,
    user_picture TEXT,
    user_first_name VARCHAR(100),
    user_last_name VARCHAR(100),
    user_gender VARCHAR(20),
    user_phone VARCHAR(50),
    user_address TEXT,
    user_license_number VARCHAR(50) NOT NULL,
    FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE,
    FOREIGN KEY (school_uuid) REFERENCES schools(school_uuid) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS driver_details;
-- +goose StatementEnd
