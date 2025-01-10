-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS students (
	student_id BIGINT PRIMARY KEY,
	student_uuid UUID UNIQUE NOT NULL,
	parent_uuid UUID NOT NULL REFERENCES users (user_uuid) ON UPDATE NO ACTION ON DELETE SET NULL,
	school_uuid UUID NOT NULL REFERENCES schools (school_uuid) ON UPDATE NO ACTION ON DELETE SET NULL,
	student_first_name VARCHAR(255) NOT NULL,
	student_last_name VARCHAR(255) NOT NULL,
	student_gender VARCHAR(20) NOT NULL,
	student_grade VARCHAR(10) NOT NULL,
	student_address TEXT NULL DEFAULT NULL,
	student_pickup_point JSON NULL DEFAULT NULL,
	created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(255) NULL DEFAULT NULL,
	updated_at TIMESTAMPTZ NULL DEFAULT NULL,
	updated_by VARCHAR(255) NULL DEFAULT NULL,
	deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
	deleted_by VARCHAR(255) NULL DEFAULT NULL
);

INSERT INTO students (student_id, student_uuid, parent_uuid, school_uuid, student_first_name, student_last_name, student_gender, student_grade, student_address, student_pickup_point, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by) VALUES
	(1735212093084389531, '33f7a0f2-be3f-4ecc-9aeb-868d7fe90aa8', '162f6603-88d9-4c7a-9571-3e0d816b3607', '5b659fa0-1d68-459f-afe6-b59e8f0e4f97', 'tono', 'sutono', 'Male', '5', '6F4V+2Q4, Seregedug Lor, Madurejo, Kec. Prambanan, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55572', '{ "latitude": -7.71159125621818, "longitude": 110.41347292301613 }', '2024-12-26 18:21:33.085246+07', 'admin2ngaglik', NULL, NULL, NULL, NULL),
	(1735211855855621869, '1aa23748-8981-4be4-bc02-611e2d07b543', '1b9be726-a2cb-44db-953b-859da29e0a96', '5b659fa0-1d68-459f-afe6-b59e8f0e4f97', 'siswa', '2ngaglik', 'Male', '5', '6F4V+2Q4, Seregedug Lor, Madurejo, Kec. Prambanan, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55572', '{ "latitude": -7.795013663547621, "longitude": 110.49443728617204 }', '2024-12-26 18:17:35.85634+07', 'admin2ngaglik', NULL, NULL, NULL, NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS students;
-- +goose StatementEnd
