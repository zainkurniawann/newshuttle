-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS vehicles (
	vehicle_id BIGINT PRIMARY KEY,
	vehicle_uuid UUID UNIQUE NOT NULL,
	school_uuid UUID NULL DEFAULT NULL,
	driver_uuid UUID NULL DEFAULT NULL,
	vehicle_name VARCHAR(50) NOT NULL,
	vehicle_number VARCHAR(20) NOT NULL,
	vehicle_type VARCHAR(20) NOT NULL,
	vehicle_color VARCHAR(20) NOT NULL,
	vehicle_seats INTEGER NOT NULL,
	vehicle_status VARCHAR(20) NULL DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(255) NULL DEFAULT NULL,
	updated_at TIMESTAMPTZ NULL DEFAULT NULL,
	updated_by VARCHAR(255) NULL DEFAULT NULL,
	deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
	deleted_by VARCHAR(255) NULL DEFAULT NULL
);


	CREATE INDEX idx_vehicle_uuid ON vehicles (vehicle_uuid);

INSERT INTO vehicles (vehicle_id, vehicle_uuid, school_uuid, driver_uuid, vehicle_name, vehicle_number, vehicle_type, vehicle_color, vehicle_seats, vehicle_status, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by) VALUES
	(1735266580364804384, 'f96ac30c-a6e6-461a-929a-a8a61efc305d', NULL, NULL, 'Honda Kharisma', 'AB 1234 CD', 'Minicar', 'Yellow', 6, 'badly damaged', '2024-12-27 09:29:40.366453+07', NULL, NULL, NULL, '2024-12-27 09:30:32.302538+07', 'admin2ngaglik'),
	(1735268008606369856, '38dfcbc0-0d10-4922-825c-a21c7154739d', NULL, NULL, 'Toyota Calya', 'AB 1234 HG', 'MPV', 'Merah', 4, 'Tersedia', '2024-12-27 09:53:28.609978+07', NULL, NULL, NULL, NULL, NULL),
	(1735269593584636268, '94744068-e344-49b4-92b1-b44b79df7a9a', NULL, NULL, 'Toyota Avanza', 'AB 1289 BG', 'MPV', 'Hitam', 5, 'Tersedia', '2024-12-27 10:19:53.585871+07', NULL, NULL, NULL, NULL, NULL),
	(1735269576513183955, 'de56332d-c919-4b0f-b604-3a947a813bc4', NULL, NULL, 'Toyota Avanza', 'AB 1156 BG', 'MPV', 'Putih', 6, 'Tersedia', '2024-12-27 10:19:36.818768+07', NULL, '2024-12-27 10:20:03.071578+07', 'sadmin', NULL, NULL),
	(1735269546158744508, '1b16f308-2095-4fb4-bbf2-06ad2d7cef82', NULL, NULL, 'Toyota Avanza', 'AB 1908 KL', 'MPV', 'Putih', 6, 'Tersedia', '2024-12-27 10:19:06.163183+07', NULL, '2024-12-27 10:20:13.121505+07', 'sadmin', NULL, NULL),
	(1735269660367769998, '8aa81b92-423e-4687-ac2f-8392264c73ed', NULL, NULL, 'Yaris', 'AB 1958 KL', 'Hatcback', 'Kuning', 4, 'Tersedia', '2024-12-27 10:21:00.368882+07', NULL, NULL, NULL, NULL, NULL),
	(1735269690492921701, '66c95c2d-aca3-48fd-8b02-b96cf88cd6e0', NULL, NULL, 'Agya', 'AB 1956 KL', 'Hatcback', 'Merah', 5, 'Tersedia', '2024-12-27 10:21:30.495489+07', NULL, NULL, NULL, NULL, NULL),
	(1735269792155846420, '9743aa6e-cd1b-4f02-8fb6-a468bde9b083', NULL, NULL, 'Toyota Vios', 'AB 1578 MN', 'Sedan', 'Hitam', 7, 'Tersedia', '2024-12-27 10:23:12.160361+07', NULL, NULL, NULL, NULL, NULL),
	(1735201380425266127, '29e68514-4e3f-402d-b118-d5c243f318ed', NULL, '5ad039d9-33e2-4331-a757-372c325fdc56', 'Koenigsegg Regera', 'BA SJABHDN A', 'Supercar', 'Yellow', 4, 'badly damaged', '2024-12-26 15:23:00.428651+07', NULL, '2024-12-27 09:28:09.612819+07', 'admin2ngaglik', NULL, NULL),
	(1735270916891366305, '9fe0a6b4-b5bc-4963-be13-5a8193c65e87', NULL, NULL, 'Ferrari 812 SuperFast 2020', 'AB 1723 HT', 'Supercar', 'Red', 2, 'Disponibile', '2024-12-27 10:41:56.89291+07', NULL, NULL, NULL, NULL, NULL),
	(1735270865045094773, '20c30b9a-a8f5-4823-9ce6-ef5520a19979', NULL, 'c143dd2b-8954-4596-93c6-3ff8bf48ac28', 'Ferrari 812 SuperFast 2020', 'AB 1349 GN', 'Supercar', 'Red', 2, 'Disponibile', '2024-12-27 10:41:05.047931+07', NULL, NULL, NULL, NULL, NULL),
	(1735270579716966458, '904c5412-f053-4f8c-8687-4d741fd6ae31', NULL, 'c77dc77e-5843-40af-a526-99d87966487f', 'SSC Tuatara', 'AB 1892 SW', 'Supercar', 'Merah', 2, 'متاح', '2024-12-27 10:36:19.7222+07', NULL, '2024-12-27 10:39:17.028388+07', 'sadmin', NULL, NULL),
	(1735269842441039545, '8cfa35c6-8aca-4db7-a75c-6a6349e84bff', NULL, '3dc5179f-4a43-4e4f-b655-12e67bc5f691', 'Wuling Alvez', 'AB 1145 DZ', 'SUV', 'Hitam', 6, 'Tersedia', '2024-12-27 10:24:02.446068+07', NULL, NULL, NULL, NULL, NULL);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vehicles;
-- +goose StatementEnd