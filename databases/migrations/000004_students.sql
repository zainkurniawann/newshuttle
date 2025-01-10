-- +goose Up
-- +goose StatementBegin
CREATE TABLE students (
    student_id BIGINT PRIMARY KEY,
    student_uuid UUID NOT NULL,
    parent_uuid UUID NOT NULL REFERENCES users(user_uuid) ON DELETE SET NULL,
    school_uuid UUID NOT NULL REFERENCES schools(school_uuid) ON DELETE CASCADE,
    student_first_name VARCHAR(255) NOT NULL,
    student_last_name VARCHAR(255) NOT NULL,
    student_gender VARCHAR(20) NOT NULL,
    student_grade VARCHAR(10) NOT NULL,
    student_address TEXT,
    student_pickup_point JSON,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_at TIMESTAMPTZ,
    updated_by VARCHAR(255),
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS students;
-- +goose StatementEnd
