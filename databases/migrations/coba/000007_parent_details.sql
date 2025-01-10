-- +goose Up
-- +goose StatementBegin
CREATE TABLE parent_details (
    user_uuid UUID PRIMARY KEY,
    user_picture TEXT,
    user_first_name VARCHAR(100),
    user_last_name VARCHAR(100),
    user_gender VARCHAR(20),
    user_phone VARCHAR(50),
    user_address TEXT,
    FOREIGN KEY (user_uuid) REFERENCES users(user_uuid) ON DELETE CASCADE
);
INSERT INTO "parent_details" ("user_uuid", "user_picture", "user_first_name", "user_last_name", "user_gender", "user_phone", "user_address") VALUES
	('1b9be726-a2cb-44db-953b-859da29e0a96', '', 'parent', 'siswa2ngaglik', 'Female', '081234567890', '6F4V+2Q4, Seregedug Lor, Madurejo, Kec. Prambanan, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55572'),
	('162f6603-88d9-4c7a-9571-3e0d816b3607', '', 'parent', 'tono', 'Female', '081234567890', '6F4V+2Q4, Seregedug Lor, Madurejo, Kec. Prambanan, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55572');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS parent_details;
-- +goose StatementEnd
