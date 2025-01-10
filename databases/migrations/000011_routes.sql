-- +goose Up
-- +goose StatementBegin
CREATE TABLE routes (
    route_id BIGINT PRIMARY KEY,
    route_uuid UUID UNIQUE NOT NULL,
    school_uuid UUID NOT NULL,
    route_name VARCHAR(100) NOT NULL,
    route_description TEXT,
    route_status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255)
);

CREATE INDEX idx_route_uuid ON routes(route_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS routes;
-- +goose StatementEnd
