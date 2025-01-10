-- +goose Up
-- +goose StatementBegin
CREATE TABLE schools (
    school_id BIGINT PRIMARY KEY,
    school_uuid UUID UNIQUE NOT NULL,
    school_name VARCHAR(255) NOT NULL,
    school_address TEXT NOT NULL,
    school_contact VARCHAR(20) NOT NULL,
    school_email VARCHAR(255) NOT NULL,
    school_description TEXT,
    school_point JSON,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255)
);

CREATE INDEX idx_school_uuid ON schools(school_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS schools;
-- +goose StatementEnd
