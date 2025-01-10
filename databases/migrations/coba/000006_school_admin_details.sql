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

INSERT INTO school_admin_details (user_uuid, school_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address) VALUES
	('39386180-9e1f-4b88-9806-182f67d2b2ef', '5b659fa0-1d68-459f-afe6-b59e8f0e4f97', '', 'Admin', 'Dua Ngaglik', 'male', '085612345678', 'Gadingan, Sinduharjo, Ngaglik, Sleman'),
	('21f9ef44-ad6b-4a2d-a906-3d969ce7eff8', 'd02e92fe-77bd-49a6-8d6e-dc2bca44fe96', '', 'Admin', 'SD Gentan', 'male', '085612345678', 'Gentan, Sinduharjo, Ngaglik, Sleman'),
	('f1b4e666-d36a-4a91-abd8-06a90edac37b', '1e35412e-57a7-4d50-985b-2f94aca055f8', '', 'Admin', 'SD 4 Pakem', 'male', '0865987651243', 'Pakem, Sleman, Yogyakarta'),
	('e7c6f54e-1e36-43e2-9861-df65db611556', 'f942ca14-26ec-41a4-b992-b4c9cfaa8760', '', 'Admin', 'SMA 2 Ngaglik', 'female', '0856918273645', 'Bandulan, Ngaglik, Sleman, Yogyakarta');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS school_admin_details;
-- +goose StatementEnd
