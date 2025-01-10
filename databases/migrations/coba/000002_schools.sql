-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS schools (
	school_id BIGINT PRIMARY KEY,
	school_uuid UUID UNIQUE NOT NULL,
	school_name VARCHAR(255) NOT NULL,
	school_address TEXT NOT NULL,
	school_contact VARCHAR(20) NOT NULL,
	school_email VARCHAR(255) NOT NULL,
	school_description TEXT NULL DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(255) NULL DEFAULT NULL,
	updated_at TIMESTAMPTZ NULL DEFAULT NULL,
	updated_by VARCHAR(255) NULL DEFAULT NULL,
	deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
	deleted_by VARCHAR(255) NULL DEFAULT NULL,
	school_point JSON NOT NULL DEFAULT '{}'
);

	CREATE INDEX idx_school_uuid ON schools (school_uuid);

INSERT INTO schools (school_id, school_uuid, school_name, school_address, school_contact, school_email, school_description, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by, school_point) VALUES
	(1735200946266228320, '1e35412e-57a7-4d50-985b-2f94aca055f8', 'SD Negeri 4 Pakem', 'Jl. Raya Turi No.23, Dero Wetan, Harjobinangun, Kec. Pakem, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55582', '027412345678', 'sekolah4pakem@gmail.com', 'sekolah 4 pakem', '2024-12-26 15:15:46.267671+07', 'sadmin', NULL, NULL, NULL, NULL, '{"latitude":-7.666575488756716,"longitude":110.41695833880696}'),
	(1735198786446709044, 'd02e92fe-77bd-49a6-8d6e-dc2bca44fe96', 'SD Negeri Gentan', 'Jl. Abimanyu No.12, Dentan, Sinduharjo, Kec. Ngaglik, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55581', '02740987654321', 'sdgentan@gmail.com', 'sekolah sd kece', '2024-12-26 14:39:46.447158+07', 'sadmin', '2024-12-26 15:20:10.542418+07', 'sadmin', NULL, NULL, '{"latitude":-7.719801577183647,"longitude":110.40610438257933}'),
	(1735198740763081999, '5b659fa0-1d68-459f-afe6-b59e8f0e4f97', 'SMP N 2 Ngaglik', 'Gadingan, Sinduharjo, Ngaglik, Sleman, Yogyakarta', '027412345678', 'sperogata@gmail.com', 'Sekolah terbaik di kota', '2024-12-26 14:39:00.765277+07', 'sadmin', '2024-12-26 15:20:28.055504+07', 'sadmin', NULL, NULL, '{"latitude":-7.716115375361339,"longitude":110.40711999603158}'),
	(1735267149608941619, 'f942ca14-26ec-41a4-b992-b4c9cfaa8760', 'SMA Negeri 2 Ngaglik', 'Jl. Besi Jangkang Km. 5, Sukoharjo, Ngaglik, Karanglo, Sukoharjo, Sleman, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55581', '0274896375121', 'sma2ngaglik@gmail.com', 'sekolah sma keren abis', '2024-12-27 09:39:09.610262+07', 'sadmin', NULL, NULL, NULL, NULL, '{"latitude":-7.705003379142222,"longitude":110.4349183102376}');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schools;
-- +goose StatementEnd