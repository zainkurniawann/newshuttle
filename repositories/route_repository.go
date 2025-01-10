package repositories

import (
	"database/sql"
	"fmt"
	"time"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RouteRepositoryInterface interface {
	CountRoutesBySchool(schoolUUID string) (int, error)
	FetchAllRoutesByAS(offset, limit int, sortField, sortDirection, schoolUUID string) ([]dto.RoutesResponseDTO, error)
	FetchSpecRouteByAS(route_name_UUID, driverUUID string) ([]entity.RouteAssignment, error)
	FetchAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)

	AddRoutes(tx *sql.Tx, route entity.Routes) (string, error)
	AddRouteAssignment(tx *sql.Tx, assignment entity.RouteAssignment) error
	
	UpdateRoute(tx *sql.Tx, route entity.Routes) error
	UpdateOrAddRouteAssignment(tx *sql.Tx, assignment entity.RouteAssignment) error

	DeleteRoute(tx *sql.Tx, routenameUUID, schoolUUID string) error
	DeleteRouteAssignments(tx *sql.Tx, routenameUUID, schoolUUID string) error

	BeginTransaction() (*sql.Tx, error)

	IsDriverAssigned(tx *sql.Tx, driverUUID string) (bool, error)
	IsStudentAssigned(tx *sql.Tx, studentUUID string) (bool, error)

	GetDriverUUIDByRouteName(routeNameUUID string) (string, error)
	ValidateDriverVehicle(driverUUID string) (bool, error)

	RouteExists(tx *sql.Tx, routenameUUID, schoolUUID string) (bool, error)
}

type routeRepository struct {
	DB *sqlx.DB
}

func NewRouteRepository(DB *sqlx.DB) *routeRepository {
	return &routeRepository{
		DB: DB,
	}
}

func (r *routeRepository) BeginTransaction() (*sql.Tx, error) {
	return r.DB.Begin()
}

func (r *routeRepository) CountRoutesBySchool(schoolUUID string) (int, error) {
	query := `
	SELECT COUNT(*)
	FROM routes
	WHERE school_uuid = $1
	`

	var total int
	err := r.DB.QueryRow(query, schoolUUID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *routeRepository) FetchAllRoutesByAS(offset, limit int, sortField, sortDirection, schoolUUID string) ([]dto.RoutesResponseDTO, error) {
	query := fmt.Sprintf(`
	SELECT 
		route_name_uuid, 
		route_name, 
		route_description, 
		created_at, 
		created_by, 
		updated_at, 
		updated_by
	FROM routes
	WHERE school_uuid = $1
	ORDER BY %s %s
	LIMIT $2 OFFSET $3
	`, sortField, sortDirection)

	rows, err := r.DB.Query(query, schoolUUID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []dto.RoutesResponseDTO

	for rows.Next() {
		var route dto.RoutesResponseDTO
		var createdAt, updatedAt sql.NullTime
		var createdBy, updatedBy sql.NullString

		err := rows.Scan(
			&route.RouteNameUUID,
			&route.RouteName,
			&route.RouteDescription,
			&createdAt,
			&createdBy,
			&updatedAt,
			&updatedBy,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			route.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		}
		if createdBy.Valid {
			route.CreatedBy = createdBy.String
		}
		if updatedAt.Valid {
			route.UpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if updatedBy.Valid {
			route.UpdatedBy = updatedBy.String
		}

		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (r *routeRepository) FetchSpecRouteByAS(routeNameUUID, driverUUID string) ([]entity.RouteAssignment, error) {
	var driverUUIDParam interface{}
	if driverUUID == "" {
		driverUUIDParam = uuid.Nil
	} else {
		driverUUIDParam = driverUUID
	}

	query := `
        SELECT 
            r.route_name_uuid,
            r.route_name,
            r.route_description,
            ra.driver_uuid,
            COALESCE(d.user_first_name, '') AS driver_first_name,
            COALESCE(d.user_last_name, '') AS driver_last_name,
            s.student_uuid,
            COALESCE(s.student_first_name, '') AS student_first_name,
            COALESCE(s.student_last_name, '') AS student_last_name,
			s.student_status,
            COALESCE(ra.student_order, 0) AS student_order
        FROM routes r
        LEFT JOIN route_assignment ra ON r.route_name_uuid = ra.route_name_uuid
        LEFT JOIN driver_details d ON ra.driver_uuid = d.user_uuid
        LEFT JOIN students s ON ra.student_uuid = s.student_uuid
        WHERE r.route_name_uuid = $1
        AND (ra.driver_uuid = $2 OR ra.driver_uuid IS NULL)
        ORDER BY ra.student_order desc
    `

	rows, err := r.DB.Query(query, routeNameUUID, driverUUIDParam)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routes: %w", err)
	}
	defer rows.Close()

	var routes []entity.RouteAssignment
	for rows.Next() {
		var route entity.RouteAssignment
		if err := rows.Scan(
			&route.RouteNameUUID,
			&route.RouteName,
			&route.RouteDescription,
			&route.DriverUUID,
			&route.DriverFirstName,
			&route.DriverLastName,
			&route.StudentUUID,
			&route.StudentFirstName,
			&route.StudentLastName,
			&route.StudentStatus,
			&route.StudentOrder,
		); err != nil {
			return nil, fmt.Errorf("failed to scan route data: %w", err)
		}

		routes = append(routes, route)
	}

	return routes, nil
}

func (repo *routeRepository) FetchAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error) {
	query := `
		SELECT
			r.route_uuid,
			r.student_uuid,
			r.driver_uuid,
			r.school_uuid,
			s.student_first_name,
			s.student_last_name,
			s.student_status,
			s.student_address,
			s.student_pickup_point,
			st.shuttle_uuid,
			st.status AS shuttle_status,
			sc.school_name,
			sc.school_point
		FROM route_assignment r
		LEFT JOIN students s ON r.student_uuid = s.student_uuid
		LEFT JOIN schools sc ON r.school_uuid = sc.school_uuid
		LEFT JOIN shuttle st ON r.student_uuid = st.student_uuid AND DATE(st.created_at) = CURRENT_DATE
		WHERE r.driver_uuid = $1 AND s.student_status = 'present'
		ORDER BY r.created_at ASC
	`
	var routes []dto.RouteResponseByDriverDTO
	err := repo.DB.Select(&routes, query, driverUUID)
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (r *routeRepository) ValidateDriverVehicle(driverUUID string) (bool, error) {
	query := `
		SELECT 
			dd.vehicle_uuid, 
			v.driver_uuid
		FROM driver_details dd
		LEFT JOIN vehicles v ON dd.user_uuid = v.driver_uuid
		WHERE dd.user_uuid = $1
	`
	var vehicleUUID sql.NullString
	var driverUUIDFromVehicle sql.NullString

	err := r.DB.QueryRow(query, driverUUID).Scan(&vehicleUUID, &driverUUIDFromVehicle)
	if err != nil {
		return false, fmt.Errorf("failed to query driver details with vehicle join: %w", err)
	}

	if !vehicleUUID.Valid || !driverUUIDFromVehicle.Valid {
		return false, nil
	}

	return true, nil
}

func (r *routeRepository) AddRoutes(tx *sql.Tx, route entity.Routes) (string, error) {
	var routeNameUUID string
	query := `
        INSERT INTO routes (
            route_id,
            route_name_uuid,
            school_uuid,
            route_name,
            route_description,
            created_at,
            created_by
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING route_name_uuid
    `

	err := tx.QueryRow(query,
		route.RouteID,
		route.RouteNameUUID,
		route.SchoolUUID,
		route.RouteName,
		route.RouteDescription,
		route.CreatedAt.Time,
		route.CreatedBy.String,
	).Scan(&routeNameUUID)
	if err != nil {
		return "", fmt.Errorf("failed to insert route: %w", err)
	}
	return routeNameUUID, nil
}

func (r *routeRepository) AddRouteAssignment(tx *sql.Tx, assignment entity.RouteAssignment) error {
	var driverCount int
	err := tx.QueryRow("SELECT COUNT(*) FROM users WHERE user_uuid = $1", assignment.DriverUUID).Scan(&driverCount)
	if err != nil {
		return fmt.Errorf("error checking driver UUID: %w", err)
	}
	if driverCount == 0 {
		return fmt.Errorf("driver not found")
	}

	var studentCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM students WHERE student_uuid = $1", assignment.StudentUUID).Scan(&studentCount)
	if err != nil {
		return fmt.Errorf("error checking student UUID: %w", err)
	}
	if studentCount == 0 {
		return fmt.Errorf("student not found")
	}

	query := `
        INSERT INTO route_assignment (
            route_id,
            route_uuid,
            driver_uuid,
            student_uuid,
            student_order,
            school_uuid,
            route_name_uuid,
            created_at,
            created_by
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err = tx.Exec(query,
		assignment.RouteID,
		assignment.RouteUUID,
		assignment.DriverUUID,
		assignment.StudentUUID,
		assignment.StudentOrder,
		assignment.SchoolUUID,
		assignment.RouteNameUUID,
		time.Now(),
		assignment.CreatedBy.String,
	)
	if err != nil {
		return fmt.Errorf("failed to insert route assignment: %w", err)
	}
	return nil
}

func (r *routeRepository) IsDriverAssigned(tx *sql.Tx, driverUUID string) (bool, error) {
	var count int
	query := `
        SELECT COUNT(*) 
        FROM route_assignment 
        WHERE driver_uuid = $1 AND deleted_at IS NULL
    `
	err := tx.QueryRow(query, driverUUID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking driver assignment: %w", err)
	}
	return count > 0, nil
}

func (r *routeRepository) IsStudentAssigned(tx *sql.Tx, studentUUID string) (bool, error) {
	var count int
	query := `
        SELECT COUNT(*) 
        FROM route_assignment 
        WHERE student_uuid = $1 AND deleted_at IS NULL
    `
	err := tx.QueryRow(query, studentUUID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking student assignment: %w", err)
	}
	return count > 0, nil
}

func (r *routeRepository) GetDriverUUIDByRouteName(routeNameUUID string) (string, error) {
	var driverUUID *string
	query := `
		SELECT 
			ra.driver_uuid
		FROM routes r
		LEFT JOIN route_assignment ra ON r.route_name_uuid = ra.route_name_uuid
		WHERE r.route_name_uuid = $1
	`
	err := r.DB.QueryRow(query, routeNameUUID).Scan(&driverUUID)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("failed to get driver UUID: %w", err)
	}
	if driverUUID == nil {
		return "", nil
	}
	return *driverUUID, nil
}

func (r *routeRepository) UpdateRoute(tx *sql.Tx, route entity.Routes) error {
	query := `
		UPDATE routes
		SET route_name = $1,
		    route_description = $2,
		    updated_at = $3,
		    updated_by = $4
		WHERE route_name_uuid = $5
	`

	_, err := tx.Exec(query,
		route.RouteName,
		route.RouteDescription,
		route.UpdatedAt.Time,
		route.UpdatedBy.String,
		route.RouteNameUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to update route: %w", err)
	}
	return nil
}

func (r *routeRepository) UpdateOrAddRouteAssignment(tx *sql.Tx, assignment entity.RouteAssignment) error {
	// Pertama, cek apakah sudah ada assignment untuk route_uuid, driver_uuid, dan student_uuid yang diberikan
	routeUUID := assignment.RouteUUID.String()  // Konversi UUID ke string
	driverUUID := assignment.DriverUUID.String()  // Konversi UUID ke string
	studentUUID := assignment.StudentUUID.String()  // Konversi UUID ke string

	exists, err := r.IsRouteAssignmentExist(tx, routeUUID, driverUUID, studentUUID)
	if err != nil {
		return fmt.Errorf("failed to check route assignment existence: %w", err)
	}

	// Jika assignment sudah ada, lakukan update
	if exists {
		query := `
			UPDATE route_assignment
			SET driver_uuid = $1,
			    student_uuid = $2,
			    student_order = $3,
			    updated_at = $4,
			    updated_by = $5
			WHERE route_uuid = $6 AND driver_uuid = $7 AND student_uuid = $8
		`

		_, err := tx.Exec(query,
			driverUUID,
			studentUUID,
			assignment.StudentOrder,
			time.Now(),
			assignment.CreatedBy.String,
			routeUUID,
			driverUUID,
			studentUUID,
		)
		if err != nil {
			return fmt.Errorf("failed to update route assignment: %w", err)
		}
	} else {
		// Jika assignment belum ada, lakukan insert
		query := `
			INSERT INTO route_assignment (route_uuid, driver_uuid, student_uuid, student_order, created_at, created_by)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

		_, err := tx.Exec(query,
			routeUUID,
			driverUUID,
			studentUUID,
			assignment.StudentOrder,
			time.Now(),
			assignment.CreatedBy.String,
		)
		if err != nil {
			return fmt.Errorf("failed to insert route assignment: %w", err)
		}
	}

	return nil
}

// Fungsi pembantu untuk memeriksa apakah route assignment sudah ada
func (r *routeRepository) IsRouteAssignmentExist(tx *sql.Tx, routeUUID, driverUUID, studentUUID string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM route_assignment 
		WHERE route_uuid = $1 AND driver_uuid = $2 AND student_uuid = $3 AND deleted_at IS NULL
	`

	err := tx.QueryRow(query, routeUUID, driverUUID, studentUUID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking route assignment existence: %w", err)
	}
	return count > 0, nil
}


func (r *routeRepository) DeleteRoute(tx *sql.Tx, routenameUUID, schoolUUID string) error {
	query := `DELETE FROM routes WHERE route_name_uuid = $1 AND school_uuid = $2`
	_, err := tx.Exec(query, routenameUUID, schoolUUID)
	if err != nil {
		return fmt.Errorf("error deleting route: %w", err)
	}
	return nil
}

func (r *routeRepository) DeleteRouteAssignments(tx *sql.Tx, routenameUUID, schoolUUID string) error {
	query := `DELETE FROM route_assignment WHERE route_name_uuid = $1 AND school_uuid = $2`
	_, err := tx.Exec(query, routenameUUID, schoolUUID)
	if err != nil {
		return fmt.Errorf("error deleting route assignments: %w", err)
	}
	return nil
}

func (r *routeRepository) RouteExists(tx *sql.Tx, routenameUUID, schoolUUID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM routes WHERE route_name_uuid = $1 AND school_uuid = $2`
	err := tx.QueryRow(query, routenameUUID, schoolUUID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking route existence: %w", err)
	}
	return count > 0, nil
}