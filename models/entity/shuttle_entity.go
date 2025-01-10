package entity

import (
	"database/sql"
	"github.com/google/uuid"
)

type Shuttle struct {
	ShuttleID    int64          `db:"shuttle_id"`
	ShuttleUUID  uuid.UUID      `db:"shuttle_uuid"`
	StudentUUID  uuid.UUID      `db:"student_uuid"`
	DriverUUID   uuid.UUID      `db:"driver_uuid"`
	Status       string         `db:"status"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	CreatedBy    sql.NullString `db:"created_by"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	UpdatedBy    sql.NullString `db:"updated_by"`
	DeletedAt    sql.NullTime   `db:"deleted_at"`
	DeletedBy    sql.NullString `db:"deleted_by"`
}
