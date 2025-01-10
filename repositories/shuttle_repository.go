package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"shuttle/models/dto"
	"shuttle/models/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShuttleRepositoryInterface interface {
	CountShuttlesByParent(parentUUID uuid.UUID) (int, error)
	FetchShuttleTrackByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error)
	FetchAllShuttleByParent(offset, limit int, sortField, sortDirection string, parentUUID uuid.UUID) ([]dto.ShuttleAllResponse, error)
	FetchAllShuttleByDriver(driverUUID uuid.UUID) ([]dto.ShuttleAllResponse, error)
	GetSpecShuttle(shuttleUUID uuid.UUID) ([]dto.ShuttleSpecResponse, error)
	SaveShuttle(shuttle entity.Shuttle) error
	UpdateShuttleStatus(shuttleUUID uuid.UUID, status string) error
}

type ShuttleRepository struct {
	DB *sqlx.DB
}

func NewShuttleRepository(DB *sqlx.DB) ShuttleRepositoryInterface {
	return &ShuttleRepository{
		DB: DB,
	}
}

func (r *ShuttleRepository) CountShuttlesByParent(parentUUID uuid.UUID) (int, error) {
    query := `
    SELECT COUNT(*)
    FROM shuttle st
    LEFT JOIN students s ON st.student_uuid = s.student_uuid
    WHERE s.parent_uuid = $1
    `

    var total int
    err := r.DB.Get(&total, query, parentUUID)
    if err != nil {
        return 0, err
    }

    return total, nil
}

func (r *ShuttleRepository) FetchShuttleTrackByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error) {
	log.Println("Executing query to fetch shuttle track for parentUUID:", parentUUID)

	query := `
		SELECT 
			st.student_uuid,
			st.shuttle_uuid,
			s.student_first_name,
			s.student_last_name,
			s.parent_uuid,
			s.school_uuid,
			sc.school_name,
			st.status AS shuttle_status,
			st.created_at,
			CURRENT_DATE AS current_date
		FROM shuttle st
		LEFT JOIN students s 
			ON s.student_uuid = st.student_uuid AND DATE(st.created_at) = CURRENT_DATE
		JOIN schools sc 
			ON s.school_uuid = sc.school_uuid
		WHERE s.parent_uuid = $1
	`
	var shuttles []dto.ShuttleResponse
	err := r.DB.Select(&shuttles, query, parentUUID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}

	log.Println("Shuttle track fetched from database:", shuttles)
	return shuttles, nil
}

func (r *ShuttleRepository) FetchAllShuttleByParent(offset, limit int, sortField, sortDirection string, parentUUID uuid.UUID) ([]dto.ShuttleAllResponse, error) {
    query := fmt.Sprintf(`
        SELECT
            st.shuttle_uuid,
            st.student_uuid,
            st.status,
            s.student_first_name,
            s.student_last_name,
            s.student_grade,
            s.student_gender,
            s.parent_uuid,
            s.school_uuid,
            sc.school_name,
            st.created_at,
            COALESCE(st.updated_at::TEXT, 'N/A') AS updated_at
        FROM shuttle st
        LEFT JOIN students s
            ON st.student_uuid = s.student_uuid
        LEFT JOIN schools sc 
            ON s.school_uuid = sc.school_uuid
        WHERE s.parent_uuid = $1
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortField, sortDirection)

    var shuttles []dto.ShuttleAllResponse
    err := r.DB.Select(&shuttles, query, parentUUID, limit, offset)
    if err != nil {
        return nil, err
    }

    return shuttles, nil
}

func (r *ShuttleRepository) FetchAllShuttleByDriver(driverUUID uuid.UUID) ([]dto.ShuttleAllResponse, error) {
	query := `
		SELECT
			st.shuttle_uuid,
			st.student_uuid,
			st.status,
			s.student_first_name,
			s.student_last_name,
			s.student_grade,
			s.student_gender,
			s.parent_uuid,
			s.school_uuid,
			sc.school_name,
			st.created_at,
			COALESCE(st.updated_at::TEXT, 'N/A') AS updated_at
		FROM shuttle st
		LEFT JOIN students s
			ON st.student_uuid = s.student_uuid
		LEFT JOIN schools sc 
			ON s.school_uuid = sc.school_uuid
		WHERE st.driver_uuid = $1 ORDER BY st.created_at DESC
	`
	var shuttles []dto.ShuttleAllResponse
	err := r.DB.Select(&shuttles, query, driverUUID)
	if err != nil {
		return nil, err
	}

	return shuttles, nil
}

func (r *ShuttleRepository) GetSpecShuttle(shuttleUUID uuid.UUID) ([]dto.ShuttleSpecResponse, error) {
	log.Println("Executing query to fetch shuttle spec for shuttleUUID:", shuttleUUID)

	query := `
		SELECT 
			st.student_uuid,
			d.user_uuid,
			d.user_username,
			dd.user_first_name AS driver_first_name,
			dd.user_last_name AS driver_last_name,
			dd.user_gender AS driver_gender,
			dd.vehicle_uuid,
			v.vehicle_name,
			v.vehicle_type,
			v.vehicle_color,
			v.vehicle_number,
			st.shuttle_uuid,
			s.student_first_name,
			s.student_last_name,
			s.student_pickup_point,
			s.parent_uuid,
			s.school_uuid,
			sc.school_name,
			sc.school_point,
			st.status AS shuttle_status,
			st.created_at,
			CURRENT_DATE AS current_date
		FROM shuttle st
		LEFT JOIN students s 
			ON s.student_uuid = st.student_uuid 
		JOIN schools sc 
			ON s.school_uuid = sc.school_uuid
		JOIN users d 
			ON st.driver_uuid = d.user_uuid
		JOIN driver_details dd 
			ON d.user_uuid = dd.user_uuid
		JOIN vehicles v 
			ON dd.vehicle_uuid = v.vehicle_uuid
		WHERE st.shuttle_uuid = $1
	`
	var shuttles []dto.ShuttleSpecResponse
	err := r.DB.Select(&shuttles, query, shuttleUUID)
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, fmt.Errorf("failed to fetch shuttle data from database: %w", err)
	}

	log.Println("Shuttle data fetched from database:", shuttles)
	for i := range shuttles {
		if err := json.Unmarshal([]byte(shuttles[i].StudentPickupPoint), &shuttles[i].StudentPickupPoint); err != nil {
			log.Println("Error unmarshalling StudentPickupPoint:", err)
		}
		if err := json.Unmarshal([]byte(shuttles[i].SchoolPoint), &shuttles[i].SchoolPoint); err != nil {
			log.Println("Error unmarshalling SchoolPoint:", err)
		}
	}

	log.Println("Successfully processed shuttle data:", shuttles)
	return shuttles, nil
}

func (r *ShuttleRepository) SaveShuttle(shuttle entity.Shuttle) error {
	// Log: Logging query execution details
	log.Printf("SaveShuttle: Preparing to execute query for shuttleID %d", shuttle.ShuttleID)

	query := `
		INSERT INTO shuttle (shuttle_id, shuttle_uuid, student_uuid, driver_uuid, status, created_at)
		VALUES (:shuttle_id, :shuttle_uuid, :student_uuid, :driver_uuid, :status, :created_at)`

	// Log: Log shuttle details before execution
	log.Printf("SaveShuttle: Shuttle details - shuttle_id: %d, shuttle_uuid: %s, student_uuid: %s, driver_uuid: %s, status: %s, created_at: %s",
		shuttle.ShuttleID, shuttle.ShuttleUUID.String(), shuttle.StudentUUID.String(), shuttle.DriverUUID.String(), shuttle.Status, shuttle.CreatedAt.Time.String())

	_, err := r.DB.NamedExec(query, shuttle)
	if err != nil {
		// Log: Error executing query
		log.Printf("SaveShuttle: Error executing query for shuttleID %d - %s", shuttle.ShuttleID, err.Error())
		return err
	}

	// Log: Successful insertion
	log.Printf("SaveShuttle: Shuttle with shuttleID %d inserted successfully", shuttle.ShuttleID)

	return nil
}

func (r *ShuttleRepository) UpdateShuttleStatus(shuttleUUID uuid.UUID, status string) error {
	query := `
		UPDATE shuttle
		SET status = :status, updated_at = NOW()
		WHERE shuttle_uuid = :shuttle_uuid`

	data := map[string]interface{}{
		"status":       status,
		"shuttle_uuid": shuttleUUID,
	}

	result, err := r.DB.NamedExec(query, data)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
