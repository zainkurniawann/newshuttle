-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
	user_id BIGINT PRIMARY KEY,
	user_uuid UUID UNIQUE NOT NULL,
	user_username VARCHAR(255) NOT NULL,
	user_email VARCHAR(255) NOT NULL,
	user_password VARCHAR(255) NOT NULL,
	user_role VARCHAR(20) NOT NULL,
	user_role_code VARCHAR(5) NULL DEFAULT NULL,
	user_status VARCHAR(20) DEFAULT 'offline',
	user_last_active TIMESTAMPTZ NULL DEFAULT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_by VARCHAR(255),
	updated_at TIMESTAMPTZ NULL DEFAULT NULL,
	updated_by VARCHAR(255) NULL DEFAULT NULL,
	deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
	deleted_by VARCHAR(255) NULL DEFAULT NULL
);

CREATE INDEX idx_user_uuid ON users(user_uuid);
CREATE INDEX idx_user_username ON users(user_username);
CREATE INDEX idx_user_email ON users(user_email);

INSERT INTO users (user_id, user_uuid, user_username, user_email, user_password, user_role, user_role_code, user_status, user_last_active, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by) VALUES
	(1735271220068874001, 'c143dd2b-8954-4596-93c6-3ff8bf48ac28', 'alessandroo', 'alessandroo@gmail.com', '$2a$10$.55ZYc7j8jXq.ZEHSiHe6uLSxXJzX/BmcjcRh1ClneREjmeiwiB4i', 'driver', 'D', 'offline', NULL, '2024-12-27 10:47:00.068535+07', 'sadmin', '2024-12-27 11:08:04.584885+07', 'sadmin', NULL, NULL),
	(1735200520903745990, '39386180-9e1f-4b88-9806-182f67d2b2ef', 'admin2ngaglik', 'admin2ngaglik@gmail.com', '$2a$10$pABBIzErdl72YKQQNuh.Bey5P9W5xpM3kZIdXTDMCIwHEHXdFrIaa', 'schooladmin', 'AS', 'offline', NULL, '2024-12-26 15:08:40.903319+07', 'sadmin', NULL, NULL, NULL, NULL),
	(1735202246346601276, '21f9ef44-ad6b-4a2d-a906-3d969ce7eff8', 'admin-sd-gentan', 'adminsdgentan@gmail.com', '$2a$10$laFJEXLBNgVRs2wgtuOkueIVBQXWrc/cPXs0rRHeB857YlP4sn0Si', 'schooladmin', 'AS', 'offline', NULL, '2024-12-26 15:37:26.346484+07', 'sadmin', NULL, NULL, NULL, NULL),
	(1735272321080861041, 'c77dc77e-5843-40af-a526-99d87966487f', 'john_speedster', 'john.sullivan.hypercar@gmail.com', '$2a$10$s9LSdn/rvlTjLcU5Ns2PbO8e5rLDgOw8IJ./uSfyBAS4SRuY7WnyO', 'driver', 'D', 'offline', NULL, '2024-12-27 11:05:21.080684+07', 'sadmin', '2024-12-27 11:27:00.1967+07', 'sadmin', NULL, NULL),
	(1735272769306215875, '3dc5179f-4a43-4e4f-b655-12e67bc5f691', 'rizky_alvez21', 'rizky.santoso.alvez@gmail.com', '$2a$10$n4bckdVE2AzLMjjGGeuSlu/CX9fAUpiLOMByj.VM/pbKHHjs2iL.u', 'driver', 'D', 'offline', NULL, '2024-12-27 11:12:49.306376+07', 'sadmin', '2024-12-27 11:27:35.804663+07', 'sadmin', NULL, NULL),
	(1735198642596447309, 'e6136e47-4079-40bf-823d-237b04172253', 'sadmin', 'superadmin@gmail.com', '$2a$10$ABcicB4I./FnB7z7rZr53uIS21cDTzfLWL/biM9LqJeEp3yVJhOyq', 'superadmin', 'SA', 'offline', '2024-12-27 11:29:40.724459+07', '2024-12-26 14:37:22.596834+07', 'jawa', '2024-12-27 10:30:15.344475+07', 'sadmin', NULL, NULL),
	(1735211855847642630, '1b9be726-a2cb-44db-953b-859da29e0a96', 'parentsiswa2ngaglik', 'parentsiswa2ngaglik@gmail.com', '$2a$10$wJBeLK9w8x25bEMv0il1E.JgIf2Iyhqs/FQHWokQNVio7jwB27E4e', 'parent', 'P', 'offline', '2024-12-27 09:35:28.081129+07', '2024-12-26 18:17:35.847643+07', 'admin2ngaglik', NULL, NULL, NULL, NULL),
	(1735267207404456327, 'e7c6f54e-1e36-43e2-9861-df65db611556', 'admin-sma-2-ngaglik', 'adminsma2ngaglik@gmail.com', '$2a$10$X0WlegxVmivEiKivZ5yCDe.l5x9PU.ljPfYCzXDZ4ZgFhY4YI3pZW', 'schooladmin', 'AS', 'offline', NULL, '2024-12-27 09:40:07.403904+07', 'sadmin', NULL, NULL, NULL, NULL),
	(1735266912823026594, 'f1b4e666-d36a-4a91-abd8-06a90edac37b', 'admin-sd-4-pakem', 'adminsd4pakem@gmail.com', '$2a$10$QPzNG8Zn0CGgfSaHJouEROikkgbUFh83nu4Pfpy78czwzHSDbxTu.', 'schooladmin', 'AS', 'offline', '2024-12-27 09:47:02.402846+07', '2024-12-27 09:35:12.823304+07', 'sadmin', NULL, NULL, NULL, NULL),
	(1735212093081560576, '162f6603-88d9-4c7a-9571-3e0d816b3607', 'parenttono', 'parenttono@gmail.com', '$2a$10$NK35UfVQFO/iVcKOC6lSeeVrS9KIsyOOeItRZ6tkto0XYmy28GHMW', 'parent', 'P', 'offline', '2024-12-27 09:47:31.823808+07', '2024-12-26 18:21:33.081582+07', 'admin2ngaglik', NULL, NULL, NULL, NULL),
	(1735211442740868417, '5ad039d9-33e2-4331-a757-372c325fdc56', 'driver-sekolah', 'driversekolah@gmail.com', '$2a$10$caGPDkSGsKqwLWeEYxbjO.zL/m/xuo/Uae1.x81uPW0m1D3rRy6jK', 'driver', 'D', 'offline', '2024-12-27 09:32:37.879502+07', '2024-12-26 18:10:42.740344+07', 'sadmin', '2024-12-27 10:35:14.37447+07', 'sadmin', NULL, NULL);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
