package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type Vehicle struct {
	ID            int64          `db:"vehicle_id"`
	UUID          uuid.UUID      `db:"vehicle_uuid"`
	SchoolUUID    *uuid.UUID     `db:"school_uuid,omitempty"`
	DriverUUID    *uuid.UUID     `db:"driver_uuid,omitempty"`
	VehicleName   string         `db:"vehicle_name"`
	VehicleNumber string         `db:"vehicle_number"`
	VehicleType   string         `db:"vehicle_type"`
	VehicleColor  string         `db:"vehicle_color"`
	VehicleSeats  int            `db:"vehicle_seats"`
	VehicleStatus string         `db:"vehicle_status"`
	CreatedAt     sql.NullTime   `db:"created_at"`
	CreatedBy     sql.NullString `db:"created_by"`
	UpdatedAt     sql.NullTime   `db:"updated_at"`
	UpdatedBy     sql.NullString `db:"updated_by"`
	DeletedAt     sql.NullTime   `db:"deleted_at"`
	DeletedBy     sql.NullString `db:"deleted_by"`
}
