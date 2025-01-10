package repositories

import (
	"fmt"
	"log"
	"shuttle/models/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StudentRepositoryInterface interface {
	CountAllStudentsWithParents(schoolUUID string) (int, error)

	FetchAllStudentsWithParents(offset int, limit int, sortField string, sortDirection string, schoolUUID string) ([]entity.Student, []entity.ParentDetails, error)
	FetchSpecStudentWithParents(studentUUID uuid.UUID, schoolUUID string) (entity.Student, entity.ParentDetails, error)
	FetchAvailableStudent(schoolUUID string) ([]entity.Student, error)
	SaveStudent(student entity.Student) error
	UpdateStudent(student entity.Student) error
	DeleteStudentWithParents(studentUUID uuid.UUID, schoolUUID, username string) error
}

type StudentRepository struct {
	db *sqlx.DB
}

func NewStudentRepository(db *sqlx.DB) StudentRepositoryInterface {
	return &StudentRepository{
		db: db,
	}
}

func (repo *StudentRepository) CountAllStudentsWithParents(schoolUUID string) (int, error) {
	var count int

	query := `SELECT COUNT(student_id) FROM students WHERE school_uuid = $1 AND deleted_at IS NULL`
	err := repo.db.Get(&count, query, schoolUUID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *StudentRepository) FetchAllStudentsWithParents(offset int, limit int, sortField string, sortDirection string, schoolUUID string) ([]entity.Student, []entity.ParentDetails, error) {
	var students []entity.Student
	var parents []entity.ParentDetails

	query := fmt.Sprintf(`
		SELECT s.student_uuid, s.parent_uuid, s.school_uuid, s.student_first_name, s.student_last_name, s.student_gender,
			s.student_grade, s.student_status, s.created_at, u.user_uuid, pd.user_first_name, pd.user_last_name, pd.user_phone, pd.user_address
		FROM students s
		INNER JOIN users u ON s.parent_uuid = u.user_uuid
		INNER JOIN parent_details pd ON s.parent_uuid = pd.user_uuid
		WHERE s.school_uuid = $1 AND u.deleted_at IS NULL AND s.deleted_at IS NULL
		ORDER BY %s %s
		LIMIT $2 OFFSET $3`,
		sortField, sortDirection)

	rows, err := repo.db.Query(query, schoolUUID, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var student entity.Student
		var parent entity.ParentDetails

		err := rows.Scan(&student.UUID, &student.ParentUUID, &student.SchoolUUID, &student.FirstName, &student.LastName,
			&student.Gender, &student.Grade, &student.Status, &student.CreatedAt, &parent.UserUUID, &parent.FirstName,
			&parent.LastName, &parent.Phone, &parent.Address)
		if err != nil {
			return nil, nil, err
		}

		students = append(students, student)
		parents = append(parents, parent)
	}

	return students, parents, nil
}

func (repo *StudentRepository) FetchSpecStudentWithParents(studentUUID uuid.UUID, schoolUUID string) (entity.Student, entity.ParentDetails, error) {
    var student entity.Student
    var parentDetails entity.ParentDetails

    query := `
    SELECT s.student_uuid, s.parent_uuid, s.school_uuid, s.student_first_name, s.student_last_name, s.student_gender,
           s.student_grade, s.student_status, s.student_address, s.student_pickup_point, s.created_at, 
           u.user_uuid, u.user_username, u.user_email, pd.user_first_name, pd.user_last_name, pd.user_phone, pd.user_address
    FROM students s
    INNER JOIN users u ON s.parent_uuid = u.user_uuid
    INNER JOIN parent_details pd ON s.parent_uuid = pd.user_uuid
    WHERE s.student_uuid = $1 AND s.school_uuid = $2 AND u.deleted_at IS NULL AND s.deleted_at IS NULL`
    
    err := repo.db.QueryRowx(query, studentUUID, schoolUUID).Scan(&student.UUID, &student.ParentUUID, &student.SchoolUUID, &student.FirstName,
        &student.LastName, &student.Gender, &student.Grade, &student.Status, &student.StudentAddress, &student.StudentPickupPoint, &student.CreatedAt,
        &parentDetails.UserUUID, &student.UserUsername, &student.UserEmail, &parentDetails.FirstName, &parentDetails.LastName, &parentDetails.Phone, &parentDetails.Address)
    if err != nil {
        return entity.Student{}, entity.ParentDetails{}, err
    }

    parentDetails = entity.ParentDetails{
        UserUUID:  parentDetails.UserUUID,
        FirstName: parentDetails.FirstName,
        LastName:  parentDetails.LastName,
        Phone:     parentDetails.Phone,
        Address:   parentDetails.Address,
    }

    return student, parentDetails, nil
}

func (repo *StudentRepository) FetchAvailableStudent(schoolUUID string) ([]entity.Student, error) {
	var students []entity.Student

	// Query SQL
	query := `
	SELECT 
		s.student_uuid,
		s.student_first_name,
		s.student_last_name
	FROM students s
	WHERE NOT EXISTS (
		SELECT 1 
		FROM route_assignment ra 
		WHERE ra.student_uuid = s.student_uuid
	) AND s.school_uuid = $1
	`

	// Scan hasil query
	rows, err := repo.db.Query(query, schoolUUID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Scan data baris per baris
	for rows.Next() {
		var student entity.Student
		err := rows.Scan(
			&student.UUID,
			&student.FirstName,
			&student.LastName,
		)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		students = append(students, student)
	}

	// Periksa apakah ada error dalam iterasi rows
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}

	return students, nil
}

func (repo *StudentRepository) SaveStudent(student entity.Student) error {
	query := `INSERT INTO students (student_id, student_uuid, parent_uuid, school_uuid, student_first_name, student_last_name,
 	student_gender, student_grade, student_status, student_address, student_pickup_point, created_by)
 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	res, err := repo.db.Exec(query, 
		student.ID, 
		student.UUID, 
		student.ParentUUID, 
		student.SchoolUUID,
		student.FirstName, 
		student.LastName, 
		student.Gender, 
		student.Grade, 
		student.Status,
		student.StudentAddress, 
		student.StudentPickupPoint.String, // Menggunakan String untuk menyimpan nilai JSON
		student.CreatedBy,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return nil
	}

	return nil
}


func (repo *StudentRepository) UpdateStudent(student entity.Student) error {
	query := `UPDATE students 
		SET student_first_name = $1, 
			student_last_name = $2, 
			student_gender = $3, 
			student_grade = $4, 
			student_address = $5, 
			student_pickup_point = $6, 
			updated_at = NOW(), 
			updated_by = $7
		WHERE student_uuid = $8 AND school_uuid = $9 AND deleted_at IS NULL`
	_, err := repo.db.Exec(query, 
		student.FirstName, 
		student.LastName, 
		student.Gender, 
		student.Grade, 
		student.StudentAddress, 
		student.StudentPickupPoint.String, // Menggunakan String untuk menyimpan nilai JSON
		student.UpdatedBy, 
		student.UUID, 
		student.SchoolUUID,
	)
	if err != nil {
		return err
	}

	return nil
}



func (repo *StudentRepository) DeleteStudentWithParents(studentUUID uuid.UUID, schoolUUID, username string) error {
	query := `UPDATE students SET deleted_at = NOW(), deleted_by = $1 WHERE student_uuid = $2 AND school_uuid = $3 AND deleted_at IS NULL`
	_, err := repo.db.Exec(query, username, studentUUID, schoolUUID)
	if err != nil {
		return err
	}

	return nil
}