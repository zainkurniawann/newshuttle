package dto

import (
	"database/sql"

	"github.com/google/uuid"
)

///////////// ROUTE ASSIGNMENT /////////////
type StudentDTO struct {
	StudentUUID      	string `json:"student_uuid"`
	StudentFirstName 	string `json:"student_first_name"`
	StudentLastName  	string `json:"student_last_name"`
	StudentStatus	  	string `json:"student_status"`
	StudentOrder 		string  `json:"student_order"`
}

type StudentReqDTO struct {
	StudentUUID 		uuid.UUID `json:"student_uuid"`
	StudentOrder 		string  `json:"student_order"`
}

type RouteAssignmentResponseDTO struct {
	DriverUUID      string         `json:"driver_uuid"`
	DriverFirstName string         `json:"driver_first_name"`
	DriverLastName  string         `json:"driver_last_name"`
	Students        []StudentDTO   `json:"students"`
}


/////////// ROUTES //////////////////////
type RoutesResponseDTO struct {
	RouteNameUUID     string                     `json:"route_name_uuid"`
	RouteName         string                     `json:"route_name"`
	RouteDescription  string                     `json:"route_description"`
	CreatedAt         string                    `json:"created_at,omitempty"`
	CreatedBy         string                    `json:"created_by,omitempty"`
	UpdatedAt         string                    `json:"updated_at,omitempty"`
	UpdatedBy         string                    `json:"updated_by,omitempty"`
	RouteAssignment   []RouteAssignmentResponseDTO `json:"route_assignment"`
}

type RoutesRequestDTO struct {
	RouteNameUUID    uuid.UUID                   `json:"route_name_uuid"`
	RouteName        string                     `json:"route_name" validate:"required"`
	RouteDescription string                     `json:"route_description" validate:"required"`
	RouteAssignment  []RouteAssignmentRequestDTO `json:"route_assignment"`
}

type RouteAssignmentRequestDTO struct {
	DriverUUID uuid.UUID   `json:"driver_uuid"`
	Students   []StudentReqDTO `json:"students"`
}

type RouteResponseByDriverDTO struct {
	RouteUUID          string         `json:"route_uuid,omitempty" db:"route_uuid"`
	StudentUUID        string         `json:"student_uuid,omitempty" db:"student_uuid"`
	DriverUUID         string         `json:"driver_uuid,omitempty" db:"driver_uuid"`
	SchoolUUID         string         `json:"school_uuid,omitempty" db:"school_uuid"`
	StudentFirstName   string         `json:"student_first_name,omitempty" db:"student_first_name"`
	StudentLastName    string        `json:"student_last_name,omitempty" db:"student_last_name"`
	StudentStatus		string			`json:"student_status,omitempty" db:"student_status"`
	StudentAddress     string         `json:"student_address,omitempty" db:"student_address"`
	StudentPickupPoint string         `json:"student_pickup_point,omitempty" db:"student_pickup_point"`
	ShuttleUUID        sql.NullString `db:"shuttle_uuid" json:"shuttle_uuid"`
	ShuttleStatus      sql.NullString `db:"shuttle_status" json:"shuttle_status"`
	SchoolName         string         `json:"school_name,omitempty" db:"school_name"`
	SchoolPoint        string         `json:"school_point,omitempty" db:"school_point"`
}