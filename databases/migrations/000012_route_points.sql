-- +goose Up
-- +goose StatementBegin
CREATE TABLE route_points (
    route_point_id BIGINT PRIMARY KEY,
    route_point_uuid UUID UNIQUE NOT NULL,
    route_uuid UUID NOT NULL REFERENCES routes(route_uuid) ON DELETE CASCADE,
    route_point_name VARCHAR(100) NOT NULL,
    route_point_order INT NOT NULL,
    route_point_latitude DOUBLE PRECISION NOT NULL,
    route_point_longitude DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255)
);

CREATE INDEX idx_route_point_uuid ON route_points(route_point_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS route_points;
-- +goose StatementEnd
