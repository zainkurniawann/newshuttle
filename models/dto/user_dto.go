package dto

import "encoding/json"

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

type UserRequestsDTO struct {
	Username  string          `json:"username" validate:"required,username,min=5,max=30"`
	Email     string          `json:"email" validate:"required,email"`
	Password  string          `json:"password" validate:"required,min=8"`
	Role      Role            `json:"role" validate:"required,role"`
	RoleCode  string          `json:"role_code"`
	Picture   string          `json:"picture"`
	FirstName string          `json:"first_name" validate:"required,max=255"`
	LastName  string          `json:"last_name" validate:"required,max=255"`
	Gender    Gender          `json:"gender" validate:"required,gender"`
	Phone     string          `json:"phone" validate:"required,phone"`
	Address   string          `json:"address" validate:"required,max=255"`
	Details   json.RawMessage `json:"details"`
}

type SchoolAdminDetailsRequestsDTO struct {
	SchoolUUID string `json:"school_uuid" validate:"required"`
}

type DriverDetailsRequestsDTO struct {
	SchoolUUID    string `json:"school_uuid"`
	VehicleUUID   string `json:"vehicle_uuid"`
	LicenseNumber string `json:"license_number" validate:"required"`
}

type UserResponseDTO struct {
	UUID       string          `json:"user_uuid"`
	Username   string          `json:"user_username"`
	Email      string          `json:"user_email"`
	Role       Role            `json:"user_role,omitempty"`
	RoleCode   string          `json:"user_role_code,omitempty"`
	Status     string          `json:"user_status"`
	LastActive string          `json:"user_last_active"`
	Details    json.RawMessage `json:"user_details"`
	CreatedAt  string          `json:"created_at,omitempty"`
	CreatedBy  string          `json:"created_by,omitempty"`
	UpdatedAt  string          `json:"updated_at,omitempty"`
	UpdatedBy  string          `json:"updated_by,omitempty"`
}

type SuperAdminDetailsResponseDTO struct {
	Picture   string `json:"user_picture,omitempty"`
	FirstName string `json:"user_first_name"`
	LastName  string `json:"user_last_name"`
	Gender    Gender `json:"user_gender"`
	Phone     string `json:"user_phone"`
	Address   string `json:"user_address,omitempty"`
}

type SchoolAdminDetailsResponseDTO struct {
	SchoolUUID string `json:"school_uuid,omitempty"`
	SchoolName string `json:"school_name"`
	Picture    string `json:"user_picture,omitempty"`
	FirstName  string `json:"user_first_name"`
	LastName   string `json:"user_last_name"`
	Gender     Gender `json:"user_gender"`
	Phone      string `json:"user_phone"`
	Address    string `json:"user_address,omitempty"`
}

type ParentDetailsResponseDTO struct {
	Picture   string `json:"user_picture,omitempty"`
	FirstName string `json:"user_first_name"`
	LastName  string `json:"user_last_name"`
	Gender    Gender `json:"user_gender"`
	Phone     string `json:"user_phone"`
	Address   string `json:"user_address,omitempty"`
}

type DriverDetailsResponseDTO struct {
	SchoolUUID    string `json:"school_uuid,omitempty"`
	SchoolName    string `json:"school_name"`
	VehicleUUID   string `json:"vehicle_uuid,omitempty"`
	VehicleNumber string `json:"vehicle_number"`
	Picture       string `json:"user_picture,omitempty"`
	FirstName     string `json:"user_first_name"`
	LastName      string `json:"user_last_name"`
	Gender        Gender `json:"user_gender"`
	Phone         string `json:"user_phone"`
	Address       string `json:"user_address,omitempty"`
	LicenseNumber string `json:"license_number"`
}
