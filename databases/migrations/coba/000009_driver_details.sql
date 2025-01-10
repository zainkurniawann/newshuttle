-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS driver_details (
	user_uuid UUID PRIMARY KEY,
	school_uuid UUID NULL DEFAULT NULL,
	vehicle_uuid UUID NULL DEFAULT NULL,
	user_picture TEXT NULL DEFAULT NULL,
	user_first_name VARCHAR(100) NULL DEFAULT NULL,
	user_last_name VARCHAR(100) NULL DEFAULT NULL,
	user_gender VARCHAR(20) NULL DEFAULT NULL,
	user_phone VARCHAR(50) NULL DEFAULT NULL,
	user_address TEXT NULL DEFAULT NULL,
	user_license_number VARCHAR(50) NOT NULL
);

INSERT INTO driver_details (user_uuid, school_uuid, vehicle_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address, user_license_number) VALUES
	('5ad039d9-33e2-4331-a757-372c325fdc56', NULL, '29e68514-4e3f-402d-b118-d5c243f318ed', '', 'Driver', 'Sekolah', 'male', '0849183748913', 'Jalan Jakal, Sleman, Yogyakarta', 'EGA308DFAN'),
	('c143dd2b-8954-4596-93c6-3ff8bf48ac28', NULL, '20c30b9a-a8f5-4823-9ce6-ef5520a19979', '', 'Alessandro', 'Ferrari', 'male', '0856198283645', 'Via Roma 24, Modena, Italia', 'G544061739250'),
	('c77dc77e-5843-40af-a526-99d87966487f', NULL, '904c5412-f053-4f8c-8687-4d741fd6ae31', '', 'John', 'Sullivan', 'male', '055548219871', '123 Hyper Drive, Beverly Hills, Los Angeles, California', 'CA-987654321'),
	('3dc5179f-4a43-4e4f-b655-12e67bc5f691', NULL, '8cfa35c6-8aca-4db7-a75c-6a6349e84bff', '', 'Rizky', 'Santoso', 'male', '085691283745', 'Jl. Raya Kemang No. 45, Bekasi, Jawa Barat', 'D-134567891');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS driver_details;
-- +goose StatementEnd