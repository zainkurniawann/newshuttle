package entity

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
)

type Role string
type Gender string

const (
	SuperAdmin  Role = "superadmin"
	SchoolAdmin Role = "schooladmin"
	Parent      Role = "parent"
	Driver      Role = "driver"

	Female Gender = "female"
	Male   Gender = "male"
)

type User struct {
	ID          int64           `db:"user_id"`
	UUID        uuid.UUID       `db:"user_uuid"`
	Username    string          `db:"user_username"`
	Email       string          `db:"user_email"`
	Password    string          `db:"user_password"`
	Role        Role            `db:"user_role"`
	RoleCode    string          `db:"user_role_code"`
	Status      string          `db:"user_status"`
	LastActive  sql.NullTime    `db:"user_last_active"`
	DetailsJSON json.RawMessage `db:"user_details"`
	CreatedAt   sql.NullTime    `db:"created_at"`
	CreatedBy   sql.NullString  `db:"created_by"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
	UpdatedBy   sql.NullString  `db:"updated_by"`
	DeletedAt   sql.NullTime    `db:"deleted_at"`
	DeletedBy   sql.NullString  `db:"deleted_by"`
}

type SuperAdminDetails struct {
	UserUUID  uuid.UUID `db:"user_uuid"`
	Picture   string    `db:"user_picture"`
	FirstName string    `db:"user_first_name"`
	LastName  string    `db:"user_last_name"`
	Gender    Gender    `db:"user_gender"`
	Phone     string    `db:"user_phone"`
	Address   string    `db:"user_address"`
}

type SchoolAdminDetails struct {
	UserUUID   uuid.UUID `db:"user_uuid"`
	SchoolUUID uuid.UUID `db:"school_uuid"`
	Picture    string    `db:"user_picture"`
	FirstName  string    `db:"user_first_name"`
	LastName   string    `db:"user_last_name"`
	Gender     Gender    `db:"user_gender"`
	Phone      string    `db:"user_phone"`
	Address    string    `db:"user_address"`
}

type ParentDetails struct {
	UserUUID  uuid.UUID `db:"user_uuid"`
	Picture   string    `db:"user_picture"`
	FirstName string    `db:"user_first_name"`
	LastName  string    `db:"user_last_name"`
	Gender    Gender    `db:"user_gender"`
	Phone     string    `db:"user_phone"`
	Address   string    `db:"user_address"`
}

type DriverDetails struct {
	UserUUID      uuid.UUID  `db:"user_uuid"`
	SchoolUUID    *uuid.UUID `db:"school_uuid"`
	VehicleUUID   *uuid.UUID `db:"vehicle_uuid"`
	Picture       string     `db:"user_picture"`
	FirstName     string     `db:"user_first_name"`
	LastName      string     `db:"user_last_name"`
	Gender        Gender     `db:"user_gender"`
	Phone         string     `db:"user_phone"`
	Address       string     `db:"user_address"`
	LicenseNumber string     `db:"user_license_number"`
}
