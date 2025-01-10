package dto

import ()

type VehicleRequestDTO struct {
	Name   string `json:"vehicle_name" validate:"required"`
	Number string `json:"vehicle_number" validate:"required"`
	Type   string `json:"vehicle_type" validate:"required"`
	Color  string `json:"vehicle_color" validate:"required"`
	Seats  int    `json:"vehicle_seats" validate:"required"`
	Status string `json:"vehicle_status" validate:"required"`
	School string `json:"school_uuid"`
}

type VehicleDriverRequestDTO struct {
	Vehicle VehicleRequestDTO `json:"vehicle"`
	Driver  UserRequestsDTO  `json:"driver"`
}

type VehicleResponseDTO struct {
	UUID       string `json:"vehicle_uuid"`
	SchoolUUID string `json:"school_uuid,omitempty"`
	SchoolName string `json:"school_name,omitempty"`
	DriverUUID string `json:"driver_uuid,omitempty"`
	DriverName string `json:"driver_name,omitempty"`
	Name       string `json:"vehicle_name"`
	Number     string `json:"vehicle_number"`
	Type       string `json:"vehicle_type"`
	Color      string `json:"vehicle_color"`
	Seats      int    `json:"vehicle_seats"`
	Status     string `json:"vehicle_status"`
	CreatedAt  string `json:"created_at,omitempty"`
	CreatedBy  string `json:"created_by,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	UpdatedBy  string `json:"updated_by,omitempty"`
}
