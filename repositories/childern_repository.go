package repositories

import (
	"shuttle/models/entity"

	"github.com/jmoiron/sqlx"
)

type ChildernRepositoryInterface interface {
	FetchAllChilderns(id string) ([]entity.Student, error)
	FetchSpecChildern(id string) (entity.Student, error)
	UpdateChildern(student entity.Student, studentUUID string) error
	UpdateChildernStatus(student entity.Student, studentUUID string) error
}

type childernRepository struct {
	DB *sqlx.DB
}

func NewChildernRepository(DB *sqlx.DB) ChildernRepositoryInterface {
	return &childernRepository{DB: DB}
}

func (repositories *childernRepository) FetchAllChilderns(id string) ([]entity.Student, error) {
	var childerns []entity.Student
	query := `
		SELECT 
			s.student_uuid,
			s.student_first_name,
			s.student_last_name,
			s.student_gender,
			s.student_grade,
			s.student_address,
			s.student_pickup_point,
			s.student_status,
			s.parent_uuid,
			s.school_uuid,
			sc.school_name
		FROM students s
		JOIN schools sc ON s.school_uuid = sc.school_uuid
		WHERE s.parent_uuid = $1
	`
	rows, err := repositories.DB.Queryx(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var childern entity.Student
		if err := rows.Scan(
			&childern.UUID,
			&childern.FirstName,
			&childern.LastName,
			&childern.Gender,
			&childern.Grade,
			&childern.StudentAddress,
			&childern.StudentPickupPoint,
			&childern.Status,
			&childern.ParentUUID,
			&childern.SchoolUUID,
			&childern.SchoolName,
		); err != nil {
			return nil, err
		}
		childerns = append(childerns, childern)
	}
	return childerns, nil
}

func (repositories *childernRepository) FetchSpecChildern(id string) (entity.Student, error) {
	var childern entity.Student
	query := `
		SELECT 
			s.student_uuid,
			s.student_first_name,
			s.student_last_name,
			s.student_gender,
			s.student_address,
			s.student_pickup_point,
			s.student_grade,
			s.student_status,
			s.parent_uuid,
			s.school_uuid,
			sc.school_name
		FROM students s
		JOIN schools sc ON s.school_uuid = sc.school_uuid
		WHERE s.student_uuid = $1
	`
	err := repositories.DB.QueryRowx(query, id).Scan(
		&childern.UUID,
		&childern.FirstName,
		&childern.LastName,
		&childern.Gender,
		&childern.StudentAddress,
		&childern.StudentPickupPoint,
		&childern.Grade,
		&childern.Status,
		&childern.ParentUUID,
		&childern.SchoolUUID,
		&childern.SchoolName,
	)
	if err != nil {
		return entity.Student{}, err
	}
	return childern, nil
}

func (repo *childernRepository) UpdateChildern(student entity.Student, studentUUID string) error {
	query := `
		UPDATE students
		SET 
			student_first_name = $1, 
			student_last_name = $2, 
			student_gender = $3, 
			student_address = $4, 
			student_pickup_point = $5,
			student_status = $6,
			updated_at = NOW(), 
			updated_by = $7
		WHERE student_uuid = $8
	`
	_, err := repo.DB.Exec(query,
		student.FirstName,
		student.LastName,
		student.Gender,
		student.StudentAddress,
		student.StudentPickupPoint,
		student.Status,
		student.UpdatedBy,
		studentUUID,
	)
	return err
}

func (repo *childernRepository) UpdateChildernStatus(student entity.Student, studentUUID string) error {
	query := `
		UPDATE students
		SET 
			student_status = $1,
			updated_at = NOW(), 
			updated_by = $2
		WHERE student_uuid = $3
	`
	_, err := repo.DB.Exec(query,
		student.Status,
		student.UpdatedBy,
		studentUUID,
	)
	return err
}