package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type School struct {
	ID          int64          `db:"school_id"`
	UUID        uuid.UUID      `db:"school_uuid"`
	Name        string         `db:"school_name"`
	Address     string         `db:"school_address"`
	Contact     string         `db:"school_contact"`
	Email       string         `db:"school_email"`
	Description string         `db:"school_description"`
	Point       sql.NullString `db:"school_point"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	CreatedBy   sql.NullString `db:"created_by"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	UpdatedBy   sql.NullString `db:"updated_by"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
	DeletedBy   sql.NullString `db:"deleted_by"`
}
