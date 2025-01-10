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

INSERT INTO super_admin_details (user_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address) VALUES
	('e6136e47-4079-40bf-823d-237b04172253', '', 'Super', 'Admin', 'male', '1297356128032', 'Jl. Sm. Km');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS super_admin_details;
-- +goose StatementEnd
