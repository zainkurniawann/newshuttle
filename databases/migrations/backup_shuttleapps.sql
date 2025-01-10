-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               PostgreSQL 16.3, compiled by Visual C++ build 1939, 64-bit
-- Server OS:                    
-- HeidiSQL Version:             12.1.0.6537
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES  */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Dumping structure for function public.delete_parent_with_no_associated_student
DELIMITER //
CREATE FUNCTION "delete_parent_with_no_associated_student"() RETURNS UNKNOWN AS $$ 
BEGIN
    -- Cek apakah orang tua (parent_uuid) masih dirujuk oleh siswa lain
    IF NOT EXISTS (
        SELECT 1 FROM students WHERE parent_uuid = OLD.parent_uuid AND deleted_at IS NULL
    ) THEN
        -- Jika tidak ada siswa lain yang merujuk, hapus orang tua (dengan mengatur deleted_at)
        UPDATE users
        SET deleted_at = NOW(), deleted_by = 'Auto Delete'
        WHERE user_uuid = OLD.parent_uuid;
    END IF;
    RETURN NULL;
END;
 $$//
DELIMITER ;

-- Dumping structure for table public.driver_details
CREATE TABLE IF NOT EXISTS "driver_details" (
	"user_uuid" UUID NOT NULL,
	"school_uuid" UUID NULL DEFAULT NULL,
	"vehicle_uuid" UUID NULL DEFAULT NULL,
	"user_picture" TEXT NULL DEFAULT NULL,
	"user_first_name" VARCHAR(100) NULL DEFAULT NULL,
	"user_last_name" VARCHAR(100) NULL DEFAULT NULL,
	"user_gender" VARCHAR(20) NULL DEFAULT NULL,
	"user_phone" VARCHAR(50) NULL DEFAULT NULL,
	"user_address" TEXT NULL DEFAULT NULL,
	"user_license_number" VARCHAR(50) NOT NULL,
	PRIMARY KEY ("user_uuid"),
	CONSTRAINT "driver_details_school_uuid_fkey" FOREIGN KEY ("school_uuid") REFERENCES "schools" ("school_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	CONSTRAINT "driver_details_user_uuid_fkey" FOREIGN KEY ("user_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT "fk_vehicle_uuid" FOREIGN KEY ("vehicle_uuid") REFERENCES "vehicles" ("vehicle_uuid") ON UPDATE NO ACTION ON DELETE SET NULL
);

-- Dumping data for table public.driver_details: 1 rows
DELETE FROM "driver_details";
/*!40000 ALTER TABLE "driver_details" DISABLE KEYS */;
INSERT INTO "driver_details" ("user_uuid", "school_uuid", "vehicle_uuid", "user_picture", "user_first_name", "user_last_name", "user_gender", "user_phone", "user_address", "user_license_number") VALUES
	('fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', NULL, 'a3935593-d1bb-46bf-8bd5-33bbe4b503e8', '', 'Kenzo', 'Lorenzo', 'male', '929715301283508', 'Jl.', '1245-2947-294657');
/*!40000 ALTER TABLE "driver_details" ENABLE KEYS */;

-- Dumping structure for table public.goose_db_version
CREATE TABLE IF NOT EXISTS "goose_db_version" (
	"id" INTEGER NOT NULL,
	"version_id" BIGINT NOT NULL,
	"is_applied" BOOLEAN NOT NULL,
	"tstamp" TIMESTAMP NOT NULL DEFAULT 'now()',
	PRIMARY KEY ("id")
);

-- Dumping data for table public.goose_db_version: 14 rows
DELETE FROM "goose_db_version";
/*!40000 ALTER TABLE "goose_db_version" DISABLE KEYS */;
INSERT INTO "goose_db_version" ("id", "version_id", "is_applied", "tstamp") VALUES
	(1, 0, 'true', '2024-12-19 11:27:48.558419'),
	(2, 1, 'true', '2024-12-19 11:27:48.579267'),
	(3, 2, 'true', '2024-12-19 11:27:48.619531'),
	(4, 3, 'true', '2024-12-19 11:27:48.66002'),
	(5, 4, 'true', '2024-12-19 11:27:48.690589'),
	(6, 5, 'true', '2024-12-19 11:27:48.708732'),
	(7, 6, 'true', '2024-12-19 11:27:48.726923'),
	(8, 7, 'true', '2024-12-19 11:27:48.743785'),
	(9, 8, 'true', '2024-12-19 11:27:48.761303'),
	(10, 9, 'true', '2024-12-19 11:27:48.774761'),
	(11, 10, 'true', '2024-12-19 11:27:48.793747'),
	(12, 11, 'true', '2024-12-19 11:27:48.801289'),
	(13, 12, 'true', '2024-12-19 11:27:48.832588'),
	(14, 20, 'true', '2024-12-19 11:27:48.866038');
/*!40000 ALTER TABLE "goose_db_version" ENABLE KEYS */;

-- Dumping structure for table public.parent_details
CREATE TABLE IF NOT EXISTS "parent_details" (
	"user_uuid" UUID NOT NULL,
	"user_picture" TEXT NULL DEFAULT NULL,
	"user_first_name" VARCHAR(100) NULL DEFAULT NULL,
	"user_last_name" VARCHAR(100) NULL DEFAULT NULL,
	"user_gender" VARCHAR(20) NULL DEFAULT NULL,
	"user_phone" VARCHAR(50) NULL DEFAULT NULL,
	"user_address" TEXT NULL DEFAULT NULL,
	PRIMARY KEY ("user_uuid"),
	CONSTRAINT "parent_details_user_uuid_fkey" FOREIGN KEY ("user_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Dumping data for table public.parent_details: 3 rows
DELETE FROM "parent_details";
/*!40000 ALTER TABLE "parent_details" DISABLE KEYS */;
INSERT INTO "parent_details" ("user_uuid", "user_picture", "user_first_name", "user_last_name", "user_gender", "user_phone", "user_address") VALUES
	('539a424e-9305-4b60-aa18-80c3e6712bb0', '', 'Jane', 'Doe', 'Female', '081234567890', '123 Main St, Springfield'),
	('57ffea3f-b9fc-408a-bb6d-fd967cfd3808', '', 'Jane', 'Doe', 'Female', '081234567890', '123 Main St, Springfield'),
	('51aabbbf-ce70-475a-ad90-79872dc86457', '', 'parent', 'budi', 'Female', '081234567890', '6F4V+2Q4, Seregedug Lor, Madurejo, Kec. Prambanan, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55572');
/*!40000 ALTER TABLE "parent_details" ENABLE KEYS */;

-- Dumping structure for table public.refresh_tokens
CREATE TABLE IF NOT EXISTS "refresh_tokens" (
	"id" BIGINT NOT NULL,
	"user_uuid" UUID NOT NULL,
	"refresh_token" TEXT NOT NULL,
	"issued_at" TIMESTAMPTZ NULL DEFAULT 'CURRENT_TIMESTAMP',
	"expired_at" TIMESTAMPTZ NOT NULL,
	"is_revoked" BOOLEAN NULL DEFAULT 'false',
	"last_used_at" TIMESTAMPTZ NULL DEFAULT NULL,
	PRIMARY KEY ("id"),
	UNIQUE INDEX "unique_user_uuid" ("user_uuid"),
	CONSTRAINT "token_fk_user" FOREIGN KEY ("user_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Dumping data for table public.refresh_tokens: 6 rows
DELETE FROM "refresh_tokens";
/*!40000 ALTER TABLE "refresh_tokens" DISABLE KEYS */;
INSERT INTO "refresh_tokens" ("id", "user_uuid", "refresh_token", "issued_at", "expired_at", "is_revoked", "last_used_at") VALUES
	(1734584561196800601, '00000000-0000-0000-0000-000000000000', '7U_lUCt5_YkxMOZ6TnNHbwd3Gpsj-pjLOW6P6Gv7EWGFu6q_B0_czh-1FOCf3zYQ3QDZRSUBwfHt-0T8evo6w8DMz_t2sAecyEHMEjQ9hSwY-PsfS5aF7xLfXX_ixCAdUNTu0K7S5SbHClv3j48nWK1Kb6xVdaGLAnA03WtJ0_7UycTwbzxyXxhq2qC1PdZUIYT7mtJTC6ByDGfN-xtfmwpjBK8R1LafWqfPdlZbmwY6YMHdugtykZM3HVJYM1t5Wdt9K_0-NqvCI8ixcCluFETEYi18L0a1pzS7txxxBW0JW9Z_xrjJipgU49h2FEhlB6rS5U0YOaIiet9X9_ZOFc-R_0L9fUs9', '2024-12-19 12:06:03.151573+07', '2025-01-03 12:06:03.15034+07', 'false', NULL),
	(1734615005720555154, 'fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', 'MxVusmhnaWVZJrMfP9bq7KR1vAM11jVz-WViBoOIZdO1h0HvPeIiI6wwIhZsR41oYwBTdmmUBKzU8r6ZQmYkbg1B-AOqOxTYXMDXvcIl4lvvBthSWfkgruj0RAt9upsXhYUYhvcB0SeCGoy8W4qOQ45IjawRszQz4U7H_c8yKXIQJv-PvS5mdY7XsHNL4XHFHfyfquJiq13O0Qw0gdm5oGRIRNDp4nnu00UDKzVIZlFRIX8F_kjKTUN7_Vo6mIeG-WpNq57rdPPhHKDm7d8ApKPnJ8GKs-xYLrbwkpnPTOAYlsZBLi0cCKldyznvHLZw-6GBAFKZPsegiGnz-5j1d3LN-a8s9F5BFrws7aJC0Ao4fF_jVsRx8KU_38MSvfc=', '2024-12-21 05:30:14.232307+07', '2025-01-05 05:30:14.230497+07', 'false', NULL),
	(1734733692863950891, '51aabbbf-ce70-475a-ad90-79872dc86457', '7BVp099IYwkKgz1w9oU3sbEAKEv5DN3NgqT4xnNWqtMnlfVCfUxc4PAH4aCsYUKB0Pp3QVIf91DW1_ZptcUtJ4RPpaIcBfZaeDLmBbvm4cV3Za5sqHQhqV7I25Vo6UEfu0ashqly0azQZzldoY76rc5AuK6HGIJ_q2Lf_5rFRuuTC4i-2Af2H-mtvDj_frUkShCtrL9EM5_PeRa-HrkQ6IunqQC2HFbSU9_zkYAwHZZYundvOkNSFkfI7sqLS7OYdAQvGYVnFcDQSd4gNiD18-3J6gmHnmtOK43Afh5vpNyLNc2OKfOP44ET05HTy-54vUgFa-Op0T2ENnZvPBcm5MUj5du1G8sYzqJYwUCdMfsuu84UNcVlXNmCy7HKrxrGiveE9ok=', '2024-12-21 05:31:11.783439+07', '2025-01-05 05:31:11.78231+07', 'false', NULL),
	(1734585217109489307, 'a7021f0d-c613-4c23-baa5-d12dd074cdb5', 'lHAuphoFNVJvFbXchlQNv91SoZl_jjfeevAHu8Sur2eozgtKBeg-0zDWSFoh1Vxpnqhsc0BSCAK7bEVZ5jgeNHqdKkUZ8oOmJkFucv1CqPkItXBTSQQToMzrcKR0evVlgXy3lDbCm5f8MPzJ0XIe7-Lxwtb_RRkCELr5llz8yOQbQCi0JfxEOkqiBN1YhuhGlJvyPm_fieysKSdu5QejnqNEswPSdeCAfn-0x7i6onBn-SQPJH1PLOLfqWFLJbhh3zth_hIZIXmhGlJdW1FJkbK3P-oIy2AhxSNx8Wf3Hi3SwuiN2K0P8YjMPI7b2DrKlS4yufnG5lV2h3TC_Dse3cMWS0dLopVUOXnuARBWBoGL6lK6L102mtRA0ZIKRh39Qw==', '2024-12-20 15:50:08.955456+07', '2025-01-04 15:50:08.954046+07', 'false', NULL),
	(1734618092018022220, '57ffea3f-b9fc-408a-bb6d-fd967cfd3808', 'TL0_v3yIP8GnsmPCB5Oh1ee82Fhpz--D6WtLd2ylneJKWipvhF5dtMUTVXb_5DjkeIy8Xfn-7Q3TcxFz0u0_kZn62QfsLiXPnc8N9BgOXbjS507NfmK7C9fiRRR5TrNgf-VhcKmrJCgOVpt0YPdtfjV4pkLDuCHFQQ7SRuSh1wZ6mqlzzxBKWh4ow7BpX1DN9-twBlubQ7AFl0EBvQBMVvmzcBUCWwtItb65X9FeZIpMG-qSRp8w0WSJcrPzfK12qTUICuPz44M-2NNkCaH2Tq7DXI_RXK3afuRo97eb1hcwDFM3Lws1f8FXxMmiqmcwEN2CqoMmswV_CstLlDjiPaOMQPIfp2idOL2VOkjcxj6_dZ5WYrvP8LNqc8vLWD_HvkiJRBiGoC1ZDoOf', '2024-12-21 05:19:37.008561+07', '2025-01-05 05:19:37.002564+07', 'false', NULL),
	(1734587810590935751, '74bc312c-576b-4bcc-b07b-e88bda1e64e6', 'Q3NxZ388mt1_1BCx892thnCT_bnJD36dIgZuRbWDzfhS8cg-tSvx44Z0ox7qe8xZk8xRZtdOA7h9_FIJupB38iPyQfBekCRn0sji98XcbpLNlYrqHHIcPGUiw-5dCfYR9RLL0bNixKYCjkv-uzhKBUGIA2OWkR50c6Zh4zx79OdBbjyeKk9ZivregbZ8Lx3NSUqh8CsdI3JfmCEBh3lCc14AAvW_gHnxHnn6gxZ9GTAO1X1UowAfTP8zcQgzH-E7sOMbShJ5nilAEBSbmc-tgD3rAeVJX79YtQkWID0FgzOAZ8odlm8fCrq27ZEihAzruNn8_WQhzm_hdjvD_KuN4aHcM3vByg8lUIb7TGX9SFPCoFwCc_7IFHk9lKDSny77FhtG7A==', '2024-12-21 05:24:22.100552+07', '2025-01-05 05:24:22.099257+07', 'false', NULL);
/*!40000 ALTER TABLE "refresh_tokens" ENABLE KEYS */;

-- Dumping structure for table public.routes
CREATE TABLE IF NOT EXISTS "routes" (
	"route_id" BIGINT NOT NULL,
	"route_uuid" UUID NOT NULL,
	"school_uuid" UUID NOT NULL,
	"route_name" VARCHAR(100) NOT NULL,
	"route_description" TEXT NULL DEFAULT NULL,
	"route_status" VARCHAR(20) NOT NULL,
	"created_at" TIMESTAMPTZ NULL DEFAULT 'CURRENT_TIMESTAMP',
	"created_by" VARCHAR(255) NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_by" VARCHAR(255) NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_by" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("route_id"),
	UNIQUE INDEX "routes_route_uuid_key" ("route_uuid"),
	INDEX "idx_route_uuid" ("route_uuid")
);

-- Dumping data for table public.routes: 0 rows
DELETE FROM "routes";
/*!40000 ALTER TABLE "routes" DISABLE KEYS */;
/*!40000 ALTER TABLE "routes" ENABLE KEYS */;

-- Dumping structure for table public.route_points
CREATE TABLE IF NOT EXISTS "route_points" (
	"route_point_id" BIGINT NOT NULL,
	"route_point_uuid" UUID NOT NULL,
	"route_uuid" UUID NOT NULL,
	"route_point_name" VARCHAR(100) NOT NULL,
	"route_point_order" INTEGER NOT NULL,
	"route_point_latitude" DOUBLE PRECISION NOT NULL,
	"route_point_longitude" DOUBLE PRECISION NOT NULL,
	"created_at" TIMESTAMPTZ NULL DEFAULT 'CURRENT_TIMESTAMP',
	"created_by" VARCHAR(255) NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_by" VARCHAR(255) NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_by" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("route_point_id"),
	UNIQUE INDEX "route_points_route_point_uuid_key" ("route_point_uuid"),
	INDEX "idx_route_point_uuid" ("route_point_uuid"),
	CONSTRAINT "route_points_route_uuid_fkey" FOREIGN KEY ("route_uuid") REFERENCES "routes" ("route_uuid") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Dumping data for table public.route_points: 0 rows
DELETE FROM "route_points";
/*!40000 ALTER TABLE "route_points" DISABLE KEYS */;
/*!40000 ALTER TABLE "route_points" ENABLE KEYS */;

-- Dumping structure for table public.schools
CREATE TABLE IF NOT EXISTS "schools" (
	"school_id" BIGINT NOT NULL,
	"school_uuid" UUID NOT NULL,
	"school_name" VARCHAR(255) NOT NULL,
	"school_address" TEXT NOT NULL,
	"school_contact" VARCHAR(20) NOT NULL,
	"school_email" VARCHAR(255) NOT NULL,
	"school_description" TEXT NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NULL DEFAULT 'CURRENT_TIMESTAMP',
	"created_by" VARCHAR(255) NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_by" VARCHAR(255) NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_by" VARCHAR(255) NULL DEFAULT NULL,
	"school_point" JSON NOT NULL DEFAULT '{}',
	PRIMARY KEY ("school_id"),
	UNIQUE INDEX "schools_school_uuid_key" ("school_uuid"),
	INDEX "idx_school_uuid" ("school_uuid")
);

-- Dumping data for table public.schools: 5 rows
DELETE FROM "schools";
/*!40000 ALTER TABLE "schools" DISABLE KEYS */;
INSERT INTO "schools" ("school_id", "school_uuid", "school_name", "school_address", "school_contact", "school_email", "school_description", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by", "school_point") VALUES
	(1734612757745200733, 'b7dc3320-e51d-4145-b939-f81d6788f153', 'Sekolah ABC', 'Jl. Pendidikan No. 1', '081234567890', 'info@sekolahabc.com', 'Sekolah terbaik di kota', '2024-12-19 19:52:37.747884+07', 'sadmin', NULL, NULL, NULL, NULL, 'null'),
	(1734612963011966717, 'a3020b18-90fe-4ca9-9e00-0e8e0654d4cd', 'Sekolah ABC', 'Jl. Pendidikan No. 1', '081234567890', 'sekolahabc@gmail.com', 'Sekolah terbaik di kota', '2024-12-19 19:56:03.013564+07', 'sadmin', NULL, NULL, NULL, NULL, 'null'),
	(1734613260993387954, '608bfb2f-4056-4d37-9257-cf7576579848', 'Sekolah ABC', 'Jl. Pendidikan No. 1', '081234567890', 'sekolahqwre@gmail.com', 'Sekolah terbaik di kota', '2024-12-19 20:01:00.99478+07', 'sadmin', NULL, NULL, NULL, NULL, '{}'),
	(1734613377111279200, '15daa126-de13-4b87-941e-54137f02fcb6', 'Sekolah ABC', 'Jl. Pendidikan No. 1', '081234567890', 'sekolahqwrere@gmail.com', 'Sekolah terbaik di kota', '2024-12-19 20:02:57.11315+07', 'sadmin', NULL, NULL, NULL, NULL, '{"latitude":-7.754455656496723,"longitude":110.38312849323236}'),
	(1734585267633713781, '281e64eb-420d-42ef-a391-e787cd08df48', 'SD Negeri 4 Pakem', 'Jl. Sm. Km', '0173613434333', 'sdn4pakem@sch.id', '', '2024-12-19 12:14:27.635782+07', 'sadmin', NULL, NULL, NULL, NULL, '{"latitude":-7.795019891634661,"longitude":110.49443642277599}');
/*!40000 ALTER TABLE "schools" ENABLE KEYS */;

-- Dumping structure for table public.school_admin_details
CREATE TABLE IF NOT EXISTS "school_admin_details" (
	"user_uuid" UUID NOT NULL,
	"school_uuid" UUID NOT NULL,
	"user_picture" TEXT NULL DEFAULT NULL,
	"user_first_name" VARCHAR(100) NULL DEFAULT NULL,
	"user_last_name" VARCHAR(100) NULL DEFAULT NULL,
	"user_gender" VARCHAR(20) NULL DEFAULT NULL,
	"user_phone" VARCHAR(50) NULL DEFAULT NULL,
	"user_address" TEXT NULL DEFAULT NULL,
	PRIMARY KEY ("user_uuid"),
	CONSTRAINT "school_admin_details_school_uuid_fkey" FOREIGN KEY ("school_uuid") REFERENCES "schools" ("school_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	CONSTRAINT "school_admin_details_user_uuid_fkey" FOREIGN KEY ("user_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Dumping data for table public.school_admin_details: 1 rows
DELETE FROM "school_admin_details";
/*!40000 ALTER TABLE "school_admin_details" DISABLE KEYS */;
INSERT INTO "school_admin_details" ("user_uuid", "school_uuid", "user_picture", "user_first_name", "user_last_name", "user_gender", "user_phone", "user_address") VALUES
	('74bc312c-576b-4bcc-b07b-e88bda1e64e6', '281e64eb-420d-42ef-a391-e787cd08df48', '', 'Evil', 'Banu', 'male', '097120832132136', 'Jl.SM.KM');
/*!40000 ALTER TABLE "school_admin_details" ENABLE KEYS */;

-- Dumping structure for table public.shuttle
CREATE TABLE IF NOT EXISTS "shuttle" (
	"shuttle_id" BIGINT NOT NULL,
	"shuttle_uuid" UUID NOT NULL DEFAULT 'gen_random_uuid()',
	"student_uuid" UUID NOT NULL,
	"driver_uuid" UUID NOT NULL,
	"status" UNKNOWN NOT NULL DEFAULT 'home',
	"created_at" TIMESTAMPTZ NULL DEFAULT 'CURRENT_TIMESTAMP',
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	PRIMARY KEY ("shuttle_id"),
	CONSTRAINT "shuttle_driver_uuid_fkey" FOREIGN KEY ("driver_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	CONSTRAINT "shuttle_student_uuid_fkey" FOREIGN KEY ("student_uuid") REFERENCES "students" ("student_uuid") ON UPDATE NO ACTION ON DELETE SET NULL
);

-- Dumping data for table public.shuttle: 5 rows
DELETE FROM "shuttle";
/*!40000 ALTER TABLE "shuttle" DISABLE KEYS */;
INSERT INTO "shuttle" ("shuttle_id", "shuttle_uuid", "student_uuid", "driver_uuid", "status", "created_at", "updated_at", "deleted_at") VALUES
	(1734615467692315567, '608882e2-450e-4b3a-93d3-acebe380917b', '27d231dd-db03-464d-8d97-ff29289e02f4', 'fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', 'going_to_school', '2024-12-19 20:37:47.692718+07', NULL, NULL),
	(1734667364103708177, 'e3e1ed4f-5377-4c49-ac00-3cb7c5334fd6', '27d231dd-db03-464d-8d97-ff29289e02f4', 'fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', 'waiting', '2024-12-19 11:02:44.103811+07', NULL, NULL),
	(1734667358589006423, '27fbd0ec-910f-4a7b-8b65-37b32c926965', '27d231dd-db03-464d-8d97-ff29289e02f4', 'fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', 'going_to_school', '2024-12-18 11:02:38.589706+07', '2024-12-20 13:36:46.692541+07', NULL),
	(1734663270057797943, 'df9ae834-649e-4f38-ae41-87f63a180076', '27d231dd-db03-464d-8d97-ff29289e02f4', 'fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', 'going_to_school', '2024-12-20 09:54:30.057545+07', NULL, NULL),
	(1734733829034850976, '92f6ed4c-4019-4ed7-8bcb-31f1d473a09d', 'd3b76cf4-f358-42f9-956a-9bf8189fd5d3', 'fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', 'waiting', '2024-12-21 05:30:29.034623+07', NULL, NULL);
/*!40000 ALTER TABLE "shuttle" ENABLE KEYS */;

-- Dumping structure for table public.students
CREATE TABLE IF NOT EXISTS "students" (
	"student_id" BIGINT NOT NULL,
	"student_uuid" UUID NOT NULL,
	"parent_uuid" UUID NOT NULL,
	"school_uuid" UUID NOT NULL,
	"student_first_name" VARCHAR(255) NOT NULL,
	"student_last_name" VARCHAR(255) NOT NULL,
	"student_gender" VARCHAR(20) NOT NULL,
	"student_grade" VARCHAR(10) NOT NULL,
	"student_address" TEXT NULL DEFAULT NULL,
	"student_pickup_point" JSON NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NULL DEFAULT 'CURRENT_TIMESTAMP',
	"created_by" VARCHAR(255) NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_by" VARCHAR(255) NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_by" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("student_id"),
	UNIQUE INDEX "unique_student_uuid" ("student_uuid"),
	CONSTRAINT "students_parent_uuid_fkey" FOREIGN KEY ("parent_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	CONSTRAINT "students_school_uuid_fkey" FOREIGN KEY ("school_uuid") REFERENCES "schools" ("school_uuid") ON UPDATE NO ACTION ON DELETE SET NULL
);

-- Dumping data for table public.students: 4 rows
DELETE FROM "students";
/*!40000 ALTER TABLE "students" DISABLE KEYS */;
INSERT INTO "students" ("student_id", "student_uuid", "parent_uuid", "school_uuid", "student_first_name", "student_last_name", "student_gender", "student_grade", "student_address", "student_pickup_point", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by") VALUES
	(1734589182544594457, 'f7426f1a-2147-4c6c-a366-9ca050942c81', '539a424e-9305-4b60-aa18-80c3e6712bb0', '281e64eb-420d-42ef-a391-e787cd08df48', 'John', 'Doe', 'Male', '5', '123 Main St, Springbed', '{"latitude":40.713876,"longitude":-74.005974}', '2024-12-19 13:19:42.546175+07', 'evilbanu', '2024-12-19 14:05:08.650114+07', 'evilbanu', '2024-12-19 14:06:28.826024+07', 'evilbanu'),
	(1734618023295754825, 'a0e3de7e-b779-4313-ae12-365d0561c199', '57ffea3f-b9fc-408a-bb6d-fd967cfd3808', '281e64eb-420d-42ef-a391-e787cd08df48', 'John', 'Doe', 'Male', '5', '123 Main St, Springbed', '{"latitude":40.713876,"longitude":-74.005974}', '2024-12-19 21:20:23.295916+07', 'evilbanu', NULL, NULL, NULL, NULL),
	(1734618021205130418, '27d231dd-db03-464d-8d97-ff29289e02f4', '57ffea3f-b9fc-408a-bb6d-fd967cfd3808', '281e64eb-420d-42ef-a391-e787cd08df48', 'John', 'Doe', 'Male', '5', '123 Main St, Springbed', '{"latitude":-7.7118030052006095,"longitude":110.41348440283387}', '2024-12-19 21:20:21.206472+07', 'evilbanu', NULL, NULL, NULL, NULL),
	(1734733654757857697, 'd3b76cf4-f358-42f9-956a-9bf8189fd5d3', '51aabbbf-ce70-475a-ad90-79872dc86457', '281e64eb-420d-42ef-a391-e787cd08df48', 'siswa', 'budi', 'Male', '5', '6F4V+2Q4, Seregedug Lor, Madurejo, Kec. Prambanan, Kabupaten Sleman, Daerah Istimewa Yogyakarta 55572', '{"latitude":7.795013663547621,"longitude":110.49443728617204}', '2024-12-21 05:27:34.758538+07', 'evilbanu', NULL, NULL, NULL, NULL);
/*!40000 ALTER TABLE "students" ENABLE KEYS */;

-- Dumping structure for table public.super_admin_details
CREATE TABLE IF NOT EXISTS "super_admin_details" (
	"user_uuid" UUID NOT NULL,
	"user_picture" TEXT NULL DEFAULT NULL,
	"user_first_name" VARCHAR(100) NULL DEFAULT NULL,
	"user_last_name" VARCHAR(100) NULL DEFAULT NULL,
	"user_gender" VARCHAR(20) NULL DEFAULT NULL,
	"user_phone" VARCHAR(50) NULL DEFAULT NULL,
	"user_address" TEXT NULL DEFAULT NULL,
	PRIMARY KEY ("user_uuid"),
	CONSTRAINT "super_admin_details_user_uuid_fkey" FOREIGN KEY ("user_uuid") REFERENCES "users" ("user_uuid") ON UPDATE NO ACTION ON DELETE CASCADE
);

-- Dumping data for table public.super_admin_details: 2 rows
DELETE FROM "super_admin_details";
/*!40000 ALTER TABLE "super_admin_details" DISABLE KEYS */;
INSERT INTO "super_admin_details" ("user_uuid", "user_picture", "user_first_name", "user_last_name", "user_gender", "user_phone", "user_address") VALUES
	('00000000-0000-0000-0000-000000000000', '', '', '', '', '', ''),
	('a7021f0d-c613-4c23-baa5-d12dd074cdb5', '', 'Super', 'Admin', 'male', '1297356128032', 'Jl. Sm. Km');
/*!40000 ALTER TABLE "super_admin_details" ENABLE KEYS */;

-- Dumping structure for table public.users
CREATE TABLE IF NOT EXISTS "users" (
	"user_id" BIGINT NOT NULL,
	"user_uuid" UUID NOT NULL,
	"user_username" VARCHAR(255) NOT NULL,
	"user_email" VARCHAR(255) NOT NULL,
	"user_password" VARCHAR(255) NOT NULL,
	"user_role" VARCHAR(20) NOT NULL,
	"user_role_code" VARCHAR(5) NULL DEFAULT NULL,
	"user_status" VARCHAR(20) NULL DEFAULT 'offline',
	"user_last_active" TIMESTAMPTZ NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NULL DEFAULT 'CURRENT_TIMESTAMP',
	"created_by" VARCHAR(255) NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_by" VARCHAR(255) NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_by" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("user_id"),
	UNIQUE INDEX "users_user_uuid_key" ("user_uuid"),
	INDEX "idx_user_uuid" ("user_uuid"),
	INDEX "idx_user_username" ("user_username"),
	INDEX "idx_user_email" ("user_email")
);

-- Dumping data for table public.users: 7 rows
DELETE FROM "users";
/*!40000 ALTER TABLE "users" DISABLE KEYS */;
INSERT INTO "users" ("user_id", "user_uuid", "user_username", "user_email", "user_password", "user_role", "user_role_code", "user_status", "user_last_active", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by") VALUES
	(1734585086352025861, 'a7021f0d-c613-4c23-baa5-d12dd074cdb5', 'sadmin', 'sadmin@gmail.com', '$2a$10$78g22p14uSRIv/nhXIG/uuKvGOM4T2UdI8qcOsRyYc3rBryq19K3a', 'superadmin', 'SA', 'offline', NULL, '2024-12-19 12:11:26.35224+07', 'faker', NULL, NULL, NULL, NULL),
	(0, '00000000-0000-0000-0000-000000000000', 'faker', 'faker@gmail.com', '$2a$10$lKS.OHbgibeu0b.hqrEQ/umVLgKZqka9gTmEo4TPb0O06vf5sQZDq', 'superadmin', 'SA', 'offline', NULL, '2024-12-19 11:29:34.700879+07', NULL, NULL, NULL, '2024-12-19 12:12:02.27207+07', 'faker'),
	(1734585340916734640, '74bc312c-576b-4bcc-b07b-e88bda1e64e6', 'evilbanu', 'banu@gmail.com', '$2a$10$600EIk4/dzkE3wwTkwtrEOqzYAkESWV6fSofIE8vcRvQEdEubJ/1u', 'schooladmin', 'AS', 'offline', NULL, '2024-12-19 12:15:40.917481+07', 'sadmin', NULL, NULL, NULL, NULL),
	(1734614990849225617, 'fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', 'knzoz', 'kenzo@gmail.com', '$2a$10$X4nIcHdeHkCq9oW1MT5EOu2pBpkoZf/D5Tq5k0zqHGRq7507VHzka', 'driver', 'D', 'offline', NULL, '2024-12-19 20:29:50.849503+07', 'sadmin', NULL, NULL, NULL, NULL),
	(1734588098028990123, '539a424e-9305-4b60-aa18-80c3e6712bb0', 'parent_john_doe', 'parent@gmail.com', '$2a$10$k3P.onZ4olkrvb2svaLSUuunVALT9ZQfGaCPL4SH22P1J96wAuG/u', 'parent', 'P', 'offline', NULL, '2024-12-19 13:01:38.02968+07', 'evilbanu', '2024-12-19 14:05:08.655013+07', 'evilbanu', '2024-12-19 14:06:28.826024+07', 'Auto Delete'),
	(1734617465032730449, '57ffea3f-b9fc-408a-bb6d-fd967cfd3808', 'parent_john_doe', 'parent@gmail.com', '$2a$10$a8PX1IXy.q8FHXM1gYL0SOTEZ5ymh25ZV3/mCa34HxkrEP7p1tTu2', 'parent', 'P', 'offline', NULL, '2024-12-19 21:11:05.03277+07', 'evilbanu', NULL, NULL, NULL, NULL),
	(1734733654747059682, '51aabbbf-ce70-475a-ad90-79872dc86457', 'parentbudi', 'parentbudi@gmail.com', '$2a$10$LCfMPlNeU8l.wjsPkB43MuVV4xpSGmj3xy6lIYY5SAW902bTgMrX2', 'parent', 'P', 'offline', NULL, '2024-12-21 05:27:34.747706+07', 'evilbanu', NULL, NULL, NULL, NULL);
/*!40000 ALTER TABLE "users" ENABLE KEYS */;

-- Dumping structure for table public.vehicles
CREATE TABLE IF NOT EXISTS "vehicles" (
	"vehicle_id" BIGINT NOT NULL,
	"vehicle_uuid" UUID NOT NULL,
	"school_uuid" UUID NULL DEFAULT NULL,
	"driver_uuid" UUID NULL DEFAULT NULL,
	"vehicle_name" VARCHAR(50) NOT NULL,
	"vehicle_number" VARCHAR(20) NOT NULL,
	"vehicle_type" VARCHAR(20) NOT NULL,
	"vehicle_color" VARCHAR(20) NOT NULL,
	"vehicle_seats" INTEGER NOT NULL,
	"vehicle_status" VARCHAR(20) NULL DEFAULT NULL,
	"created_at" TIMESTAMPTZ NULL DEFAULT 'CURRENT_TIMESTAMP',
	"created_by" VARCHAR(255) NULL DEFAULT NULL,
	"updated_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"updated_by" VARCHAR(255) NULL DEFAULT NULL,
	"deleted_at" TIMESTAMPTZ NULL DEFAULT NULL,
	"deleted_by" VARCHAR(255) NULL DEFAULT NULL,
	PRIMARY KEY ("vehicle_id"),
	UNIQUE INDEX "vehicles_vehicle_uuid_key" ("vehicle_uuid"),
	INDEX "idx_vehicle_uuid" ("vehicle_uuid"),
	CONSTRAINT "fk_driver_uuid" FOREIGN KEY ("driver_uuid") REFERENCES "driver_details" ("user_uuid") ON UPDATE NO ACTION ON DELETE SET NULL,
	CONSTRAINT "vehicles_school_uuid_fkey" FOREIGN KEY ("school_uuid") REFERENCES "schools" ("school_uuid") ON UPDATE NO ACTION ON DELETE SET NULL
);

-- Dumping data for table public.vehicles: 1 rows
DELETE FROM "vehicles";
/*!40000 ALTER TABLE "vehicles" DISABLE KEYS */;
INSERT INTO "vehicles" ("vehicle_id", "vehicle_uuid", "school_uuid", "driver_uuid", "vehicle_name", "vehicle_number", "vehicle_type", "vehicle_color", "vehicle_seats", "vehicle_status", "created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by") VALUES
	(1734614950094675387, 'a3935593-d1bb-46bf-8bd5-33bbe4b503e8', NULL, 'fa835e31-8ea4-44d7-b98d-bcf8b2f6c92b', 'Koenigsegg Regera', 'BA SJABHDN A', 'Supercar', 'Yellow', 4, 'Need Repair', '2024-12-19 20:29:10.101473+07', NULL, NULL, NULL, NULL, NULL);
/*!40000 ALTER TABLE "vehicles" ENABLE KEYS */;

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
