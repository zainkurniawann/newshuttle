package entity

import (
	"database/sql"

	"github.com/google/uuid"
)

// type Student struct {
// 	StudentID        int64          `db:"student_id"`
// 	StudentUUID      uuid.UUID      `db:"student_uuid"`
// 	ParentUUID       uuid.UUID      `db:"parent_uuid"`
// 	SchoolUUID       uuid.UUID      `db:"school_uuid"`
// 	StudentFirstName string         `db:"student_first_name"`
// 	StudentLastName  string         `db:"student_last_name"`
// 	StudentGender    string         `db:"student_gender"`
// 	StudentGrade     string         `db:"student_grade"`
// 	StudentAddress   sql.NullString `db:"student_address"` // Menambahkan field student_address
// 	StudentPickupPoint sql.NullString `db:"student_pickup_point"` // Menambahkan field student_pickup_point
// 	ShuttleStatus sql.NullString // Menambahkan field student_pickup_point
// 	CreatedAt        sql.NullTime   `db:"created_at"`
// 	CreatedBy        sql.NullString `db:"created_by"`
// 	UpdatedAt        sql.NullTime   `db:"updated_at"`
// 	UpdatedBy        sql.NullString `db:"updated_by"`
// 	DeletedAt        sql.NullTime   `db:"deleted_at"`
// 	DeletedBy        sql.NullString `db:"deleted_by"`
// }

type Student struct {
	ID        int64          `db:"student_id"`
	UUID      uuid.UUID      `db:"student_uuid"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
	Grade     string         `db:"student_grade"`
	StudentAddress   sql.NullString `db:"student_address"` // Menambahkan field student_address
	StudentPickupPoint sql.NullString `db:"student_pickup_point"`
	Gender     string         `db:"student_gender"`
	Status 		string 	`db:"student_status"`
	ParentID  sql.NullInt64  `db:"parent_id"`
	ParentUUID sql.NullString `db:"parent_uuid"`
	SchoolID  int64          `db:"school_id"`
	SchoolUUID uuid.UUID     `db:"school_uuid"`
	UserUsername	string
	UserEmail	string
    SchoolName string   
	ShuttleStatus sql.NullString  
	CreatedAt sql.NullTime   `db:"created_at"`
	CreatedBy sql.NullString `db:"created_by"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

