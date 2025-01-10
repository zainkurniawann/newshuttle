-- +goose Up
-- +goose StatementBegin
-- Create ENUM type for shuttle status
CREATE TYPE shuttle_status AS ENUM (
    'home',
    'waiting_to_be_taken_to_school',
    'going_to_school',
    'at_school',
    'waiting_to_be_taken_to_home',
    'going_to_home'
);

ALTER TYPE shuttle_status
    OWNER TO postgres;

-- Create the shuttle table with shuttle_status ENUM type for the status column
CREATE TABLE IF NOT EXISTS shuttle (
    shuttle_id BIGINT PRIMARY KEY,
    shuttle_uuid UUID NOT NULL,
    student_uuid UUID NOT NULL,
    driver_uuid UUID NOT NULL,
    status shuttle_status NOT NULL DEFAULT 'home',
    created_at TIMESTAMPTZ NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NULL DEFAULT NULL,
    deleted_at TIMESTAMPTZ NULL DEFAULT NULL,
    FOREIGN KEY (driver_uuid) REFERENCES users (user_uuid) ON UPDATE NO ACTION ON DELETE SET NULL,
    FOREIGN KEY (student_uuid) REFERENCES students (student_uuid) ON UPDATE NO ACTION ON DELETE SET NULL
);

-- Insert sample data
INSERT INTO shuttle (shuttle_id, shuttle_uuid, student_uuid, driver_uuid, created_at, updated_at, deleted_at, status) VALUES
    (1735218943426639324, 'e2a044c8-851d-49c6-b18f-f5eb8a97d4ff', '1aa23748-8981-4be4-bc02-611e2d07b543', '5ad039d9-33e2-4331-a757-372c325fdc56', '2024-12-26 20:15:43.426792+07', '2024-12-26 20:18:03.81393+07', NULL, 'home'),
    (1735308180495317259, '49c708e6-c3d8-407a-87c6-f16a208f0a2b', '33f7a0f2-be3f-4ecc-9aeb-868d7fe90aa8', '5ad039d9-33e2-4331-a757-372c325fdc56', '2024-12-27 21:03:00.495055+07', NULL, NULL, 'waiting_to_be_taken_to_school'),
    (1735219097654021499, '56be1f67-5731-4680-8c22-d2806fa9b8e6', '33f7a0f2-be3f-4ecc-9aeb-868d7fe90aa8', '5ad039d9-33e2-4331-a757-372c325fdc56', '2024-12-26 20:18:17.654631+07', '2024-12-26 21:06:44.683714+07', NULL, 'going_to_school'),
    (1735308178039420661, '04199185-7931-43ae-878c-2dc6e52fe4fc', '1aa23748-8981-4be4-bc02-611e2d07b543', '5ad039d9-33e2-4331-a757-372c325fdc56', '2024-12-27 21:02:58.039166+07', '2024-12-27 14:01:36.809669+07', NULL, 'at_school');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shuttle;
DROP TYPE IF EXISTS shuttle_status;
-- +goose StatementEnd
