package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

type RouteAssignment struct {
	RouteID          int64          `db:"route_id"`
	RouteUUID        uuid.UUID      `db:"route_uuid"`
	RouteNameUUID	string
	DriverUUID		 uuid.UUID		`db:"driver_uuid"`
    DriverFirstName   string
    DriverLastName    string
	UserUsername 	string			`json:"user_username"`
	StudentUUID		 uuid.UUID		`db:"Student_uuid"`
	StudentFirstName  string
    StudentLastName   string
	StudentStatus	string
	StudentOrder	string
	StudentName		string			`json:"student_name"`
	SchoolUUID       uuid.UUID      `db:"school_uuid"`
	RouteName        string         `db:"route_name"`
	RouteDescription string         `db:"route_description"`
	CreatedAt        sql.NullTime   `db:"created_at"`
	CreatedBy        sql.NullString `db:"created_by"`
	UpdatedAt        sql.NullTime   `db:"updated_at"`
	UpdatedBy        sql.NullString `db:"updated_by"`
	DeletedAt        sql.NullTime   `db:"deleted_at"`
	DeletedBy        sql.NullString `db:"deleted_by"`
}

type Routes struct {
	RouteID          		int64          `db:"route_id"`
	RouteNameUUID           uuid.UUID      `db:"route_name_uuid"`
	SchoolUUID       		uuid.UUID      `db:"school_uuid"`
	RouteName      			string         `db:"route_name"`
	RouteDescription     	string         `db:"route_description"`
	CreatedAt          		sql.NullTime   `db:"created_at"`
	CreatedBy           	sql.NullString `db:"created_by"`
	UpdatedAt           	sql.NullTime   `db:"updated_at"`
	UpdatedBy           	sql.NullString `db:"updated_by"`
	DeletedAt           	sql.NullTime   `db:"deleted_at"`
	DeletedBy           	sql.NullString `db:"deleted_by"`
}