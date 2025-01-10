-- +goose Up
-- +goose StatementBegin
ALTER TABLE vehicles
    ADD CONSTRAINT fk_driver_uuid FOREIGN KEY (driver_uuid) REFERENCES driver_details(user_uuid) ON DELETE SET NULL;

ALTER TABLE driver_details
    ADD CONSTRAINT fk_vehicle_uuid FOREIGN KEY (vehicle_uuid) REFERENCES vehicles(vehicle_uuid) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE vehicles
    DROP CONSTRAINT fk_driver_uuid;

ALTER TABLE driver_details
    DROP CONSTRAINT fk_vehicle_uuid;
-- +goose StatementEnd
