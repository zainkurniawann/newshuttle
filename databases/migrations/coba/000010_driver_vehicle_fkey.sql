-- +goose Up
-- +goose StatementBegin
ALTER TABLE vehicles 
    ADD CONSTRAINT fk_driver_uuid_1 FOREIGN KEY (driver_uuid) REFERENCES driver_details (user_uuid) ON UPDATE NO ACTION ON DELETE NO ACTION,
    ADD CONSTRAINT fk_driver_uuid_2 FOREIGN KEY (driver_uuid) REFERENCES users (user_uuid) ON UPDATE NO ACTION ON DELETE SET NULL,
    ADD CONSTRAINT fk_school_uuid FOREIGN KEY (school_uuid) REFERENCES schools (school_uuid) ON UPDATE NO ACTION ON DELETE SET NULL;

ALTER TABLE driver_details
    ADD CONSTRAINT fk_school_uuid FOREIGN KEY (school_uuid) REFERENCES schools (school_uuid) ON UPDATE NO ACTION ON DELETE SET NULL,
    ADD CONSTRAINT fk_user_uuid FOREIGN KEY (user_uuid) REFERENCES users (user_uuid) ON UPDATE NO ACTION ON DELETE CASCADE,
    ADD CONSTRAINT fk_vehicle_uuid FOREIGN KEY (vehicle_uuid) REFERENCES vehicles (vehicle_uuid) ON UPDATE NO ACTION ON DELETE SET NULL;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE vehicles
    DROP CONSTRAINT IF EXISTS fk_driver_uuid_1,
    DROP CONSTRAINT IF EXISTS fk_driver_uuid_2,
    DROP CONSTRAINT IF EXISTS fk_school_uuid;

ALTER TABLE driver_details
    DROP CONSTRAINT IF EXISTS fk_school_uuid,
    DROP CONSTRAINT IF EXISTS fk_user_uuid,
    DROP CONSTRAINT IF EXISTS fk_vehicle_uuid;
-- +goose StatementEnd
