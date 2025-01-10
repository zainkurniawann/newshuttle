-- +goose Up
-- +goose StatementBegin
CREATE TABLE vehicles (
    vehicle_id BIGINT PRIMARY KEY,
    vehicle_uuid UUID UNIQUE NOT NULL,
    school_uuid UUID NULL REFERENCES schools(school_uuid) ON DELETE SET NULL,
    driver_uuid UUID NULL,
    vehicle_name VARCHAR(50) NOT NULL,
    vehicle_number VARCHAR(20) NOT NULL,
    vehicle_type VARCHAR(20) NOT NULL,
    vehicle_color VARCHAR(20) NOT NULL,
    vehicle_seats INT NOT NULL,
    vehicle_status VARCHAR(20),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255)
);

CREATE INDEX idx_vehicle_uuid ON vehicles(vehicle_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vehicles;
-- +goose StatementEnd
