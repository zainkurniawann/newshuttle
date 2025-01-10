package repositories

import (
	"database/sql"
	"fmt"
	"shuttle/models/entity"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SchoolRepositoryInterface interface {
	FetchAllSchools(offset, limit int, sortField, sortDirection string) ([]entity.School, map[string][]entity.SchoolAdminDetails, error)
	FetchSpecSchool(uuid string) (entity.School, []entity.SchoolAdminDetails, error)
	SaveSchool(entity.School) error
	UpdateSchool(entity.School) error
	DeleteSchool(entity.School) error
	CountSchools() (int, error)
}

type schoolRepository struct {
	DB *sqlx.DB
}

func NewSchoolRepository(DB *sqlx.DB) SchoolRepositoryInterface {
	return &schoolRepository{
		DB: DB,
	}
}

func (repositories *schoolRepository) FetchAllSchools(offset, limit int, sortField, sortDirection string) ([]entity.School, map[string][]entity.SchoolAdminDetails, error) {
	var schools []entity.School
	var adminMap = make(map[string][]entity.SchoolAdminDetails)

	query := fmt.Sprintf(`
        SELECT 
			s.school_uuid, 
			s.school_name, 
			s.school_address, 
			s.school_contact, 
			s.school_email, 
			s.created_at,
			COALESCE(
				STRING_AGG(
					CASE
						WHEN u.deleted_at IS NULL THEN sad.user_uuid::TEXT
						ELSE NULL
					END, ', '
				),
				NULL
			) AS user_uuids,
			COALESCE(
				STRING_AGG(
					CASE
						WHEN u.deleted_at IS NULL THEN sad.school_uuid::TEXT
						ELSE NULL
					END, ', '
				),
				NULL
			) AS admin_school_uuids,
			COALESCE(
				STRING_AGG(
					CASE
						WHEN u.deleted_at IS NULL THEN sad.user_first_name
						ELSE 'N/A'
					END, ', '
				),
				'N/A'
			) AS user_first_names,
			COALESCE(
				STRING_AGG(
					CASE
						WHEN u.deleted_at IS NULL THEN sad.user_last_name
						ELSE 'N/A'
					END, ', '
				),
				'N/A'
			) AS user_last_names
		FROM schools s
		LEFT JOIN school_admin_details sad ON s.school_uuid = sad.school_uuid
		LEFT JOIN users u ON sad.user_uuid = u.user_uuid
		WHERE s.deleted_at IS NULL
		GROUP BY
			s.school_id,
			s.school_uuid, 
			s.school_name, 
			s.school_address, 
			s.school_contact, 
			s.school_email, 
			s.created_at
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sortField, sortDirection)

	rows, err := repositories.DB.Queryx(query, limit, offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var school entity.School
		var userUUIDs, adminSchoolUUIDs, firstNames, lastNames sql.NullString // Use NullString to handle NULL values
	
		if err := rows.Scan(&school.UUID, &school.Name, &school.Address, &school.Contact, &school.Email, &school.CreatedAt,
			&userUUIDs, &adminSchoolUUIDs, &firstNames, &lastNames); err != nil {
			return nil, nil, err
		}

		schools = append(schools, school)

		// Handle the case when STRING_AGG result is NULL (which means no values)
		if userUUIDs.Valid && userUUIDs.String != "" {
			// Proses hasil `STRING_AGG` untuk memisahkan UUID individual
			userUUIDList := strings.Split(userUUIDs.String, ", ")
			adminSchoolUUIDList := strings.Split(adminSchoolUUIDs.String, ", ")
			firstNameList := strings.Split(firstNames.String, ", ")
			lastNameList := strings.Split(lastNames.String, ", ")

			// Gabungkan menjadi detail admin
			for i := range userUUIDList {
				admin := entity.SchoolAdminDetails{
					UserUUID:   uuid.MustParse(userUUIDList[i]),
					SchoolUUID: uuid.MustParse(adminSchoolUUIDList[i]),
					FirstName:  firstNameList[i],
					LastName:   lastNameList[i],
				}
				adminMap[school.UUID.String()] = append(adminMap[school.UUID.String()], admin)
			}
		} else {
			// If the userUUID is NULL, handle appropriately (e.g., no admins)
			adminMap[school.UUID.String()] = []entity.SchoolAdminDetails{}
		}
	}

	return schools, adminMap, nil
}

func (repositories *schoolRepository) FetchSpecSchool(id string) (entity.School, []entity.SchoolAdminDetails, error) {
	var school entity.School
	var admin []entity.SchoolAdminDetails
	var userUUIDs, adminSchoolUUIDs, firstNames, lastNames sql.NullString // Use sql.NullString for nullable fields

	query := `
		SELECT s.school_uuid, s.school_name, s.school_address, s.school_contact, s.school_email, s.school_description, s.school_point, s.created_at,
			s.created_by, s.updated_at, s.updated_by, 
			COALESCE(
				STRING_AGG(
					CASE
						WHEN u.deleted_at IS NULL THEN sad.user_uuid::TEXT
						ELSE NULL
					END, ', '
				),
				NULL
			) AS user_uuids,
			COALESCE(
				STRING_AGG(
					CASE
						WHEN u.deleted_at IS NULL THEN sad.school_uuid::TEXT
						ELSE NULL
					END, ', '
				),
				NULL
			) AS admin_school_uuids,
			COALESCE(
				STRING_AGG(
					CASE
						WHEN u.deleted_at IS NULL THEN sad.user_first_name
						ELSE 'N/A'
					END, ', '
				),
				'N/A'
			) AS user_first_names,
			COALESCE(
				STRING_AGG(
					CASE
						WHEN u.deleted_at IS NULL THEN sad.user_last_name
						ELSE 'N/A'
					END, ', '
				),
				'N/A'
			) AS user_last_names
		FROM schools s
		LEFT JOIN school_admin_details sad ON s.school_uuid = sad.school_uuid
		LEFT JOIN users u ON sad.user_uuid = u.user_uuid
		WHERE s.deleted_at IS NULL AND s.school_uuid = $1
		GROUP BY
			s.school_id,
			s.school_uuid, 
			s.school_name, 
			s.school_address, 
			s.school_contact, 
			s.school_email, 
			s.created_at
	`

	err := repositories.DB.QueryRowx(query, id).Scan(
		&school.UUID, &school.Name, &school.Address, &school.Contact, &school.Email, &school.Description, &school.Point, &school.CreatedAt,
		&school.CreatedBy, &school.UpdatedAt, &school.UpdatedBy, &userUUIDs, &adminSchoolUUIDs, &firstNames, &lastNames,
	)
	if err != nil {
		return entity.School{}, []entity.SchoolAdminDetails{}, err
	}

	// Check if the STRING_AGG values are valid and non-empty
	if userUUIDs.Valid && userUUIDs.String != "" {
		// Process STRING_AGG results into lists
		userUUIDList := strings.Split(userUUIDs.String, ", ")
		adminSchoolUUIDList := strings.Split(adminSchoolUUIDs.String, ", ")
		firstNameList := strings.Split(firstNames.String, ", ")
		lastNameList := strings.Split(lastNames.String, ", ")

		// Add the first admin detail (if available) to the admin struct
		for i := range userUUIDList {
			admin = append(admin, entity.SchoolAdminDetails{
				UserUUID:   uuid.MustParse(userUUIDList[i]),
				SchoolUUID: uuid.MustParse(adminSchoolUUIDList[i]),
				FirstName:  firstNameList[i],
				LastName:   lastNameList[i],
			})
		}
	}

	return school, admin, nil
}

func (r *schoolRepository) SaveSchool(school entity.School) error {
	// Memeriksa apakah Point valid
	var point interface{}
	if school.Point.Valid {
		point = school.Point.String
	} else {
		point = nil // Menetapkan NULL untuk school_point jika tidak valid
	}

	query := `INSERT INTO schools (school_id, school_uuid, school_name, school_address, school_contact, school_email, school_description, school_point, created_by)
			  VALUES (:school_id, :school_uuid, :school_name, :school_address, :school_contact, :school_email, :school_description, :school_point, :created_by)`
	
	// Menyertakan value untuk school_point
	_, err := r.DB.NamedExec(query, map[string]interface{}{
		"school_id":        school.ID,
		"school_uuid":      school.UUID,
		"school_name":      school.Name,
		"school_address":   school.Address,
		"school_contact":   school.Contact,
		"school_email":     school.Email,
		"school_description": school.Description,
		"school_point":     point, // Menyimpan NULL jika tidak valid
		"created_by":       school.CreatedBy,
	})
	if err != nil {
		return err
	}

	return nil
}


func (r *schoolRepository) UpdateSchool(school entity.School) error {
	query := `
		UPDATE schools SET school_name = :school_name, school_address = :school_address, school_contact = :school_contact, school_email = :school_email, school_description = :school_description, school_point = :school_point, updated_at = :updated_at, updated_by = :updated_by 
		WHERE school_uuid = :school_uuid`
	_, err := r.DB.NamedExec(query, school)
	if err != nil {
		return err
	}

	return nil
}

func (r *schoolRepository) DeleteSchool(school entity.School) error {
	query := `UPDATE schools SET deleted_at = :deleted_at, deleted_by = :deleted_by WHERE school_uuid = :school_uuid`
	_, err := r.DB.NamedExec(query, school)
	if err != nil {
		return err
	}

	return nil
}

func (repositories *schoolRepository) CountSchools() (int, error) {
	var total int

	query := `
        SELECT COUNT(*) 
        FROM schools 
		WHERE deleted_at IS NULL
    `

	if err := repositories.DB.Get(&total, query); err != nil {
		return 0, err
	}

	return total, nil
}
