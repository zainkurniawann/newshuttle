package dto

import "database/sql"

type ShuttleRequest struct {
	StudentUUID string `json:"student_uuid" validate:"required,uuid4"`
	Status      string `json:"status" validate:"required"`
}

type ShuttleResponse struct {
	StudentUUID     string `db:"student_uuid" json:"student_uuid"`
	ShuttleUUID     string `db:"shuttle_uuid" json:"shuttle_uuid"`
	StudentFirstName string `db:"student_first_name" json:"student_first_name"`
	StudentLastName string `db:"student_last_name" json:"student_last_name"`
	ParentUUID      string `db:"parent_uuid" json:"parent_uuid"`
	SchoolUUID      string `db:"school_uuid" json:"school_uuid"`
	SchoolName      string `db:"school_name" json:"school_name"`
	ShuttleStatus   string `db:"shuttle_status" json:"shuttle_status"`
	CreatedAt       string `db:"created_at" json:"created_at"`
	CurrentDate     string `db:"current_date" json:"current_date"`
}

type ShuttleAllResponse struct {
	ShuttleUUID     string `db:"shuttle_uuid" json:"shuttle_uuid"`
	StudentUUID      string         `db:"student_uuid" json:"student_uuid"`
	Status           string         `db:"status" json:"status"`
	StudentFirstName string         `db:"student_first_name" json:"student_first_name"`
	StudentLastName  string         `db:"student_last_name" json:"student_last_name"`
	StudentGrade     string         `db:"student_grade" json:"student_grade"`
	StudentGender    string         `db:"student_gender" json:"student_gender"`
	ParentUUID       string         `db:"parent_uuid" json:"parent_uuid"`
	SchoolUUID       string         `db:"school_uuid" json:"school_uuid"`
	SchoolName       string         `db:"school_name" json:"school_name"`
	CreatedAt        string         `db:"created_at" json:"created_at"`
	UpdatedAt        sql.NullString `db:"updated_at" json:"updated_at"`
}

type ShuttleSpecResponse struct {
	StudentUUID        string `db:"student_uuid" json:"student_uuid"`
	DriverUUID         string `db:"user_uuid" json:"driver_uuid"`
	DriverUsername     string `db:"user_username" json:"driver_username"`
	DriverFirstName    string `db:"driver_first_name" json:"driver_first_name"`
	DriverLastName     string `db:"driver_last_name" json:"driver_last_name"`
	DriverGender       string `db:"driver_gender" json:"driver_gender"`
	VehicleUUID        string `db:"vehicle_uuid" json:"vehicle_uuid"`
	VehicleName        string `db:"vehicle_name" json:"vehicle_name"`
	VehicleType        string `db:"vehicle_type" json:"vehicle_type"`
	VehicleColor       string `db:"vehicle_color" json:"vehicle_color"`
	VehicleNumber      string `db:"vehicle_number" json:"vehicle_number"`
	ShuttleUUID        string `db:"shuttle_uuid" json:"shuttle_uuid"`
	StudentFirstName   string `db:"student_first_name" json:"student_first_name"`
	StudentLastName    string `db:"student_last_name" json:"student_last_name"`
	StudentPickupPoint string `db:"student_pickup_point" json:"student_pickup_point"` // Format JSON
	ParentUUID         string `db:"parent_uuid" json:"parent_uuid"`
	SchoolUUID         string `db:"school_uuid" json:"school_uuid"`
	SchoolName         string `db:"school_name" json:"school_name"`
	SchoolPoint        string `db:"school_point" json:"school_point"` // Format JSON
	ShuttleStatus      string `db:"shuttle_status" json:"shuttle_status"`
	CreatedAt          string `db:"created_at" json:"created_at"`
	CurrentDate        string `db:"current_date" json:"current_date"`
}
