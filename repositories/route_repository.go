package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RouteRepositoryInterface interface {
	CountRoutesBySchool(schoolUUID string) (int, error)

	FetchAllRoutesByAS(offset, limit int, sortField, sortDirection, schoolUUID string) ([]dto.RoutesResponseDTO, error)
	FetchSpecRouteByAS(routeNameUUID, driverUUID string) ([]entity.RouteAssignment, error)
	FetchAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)

	AddRoutes(tx *sql.Tx, route entity.Routes) (string, error)
	AddRouteAssignment(tx *sql.Tx, assignment entity.RouteAssignment) error

	UpdateRouteDetails(route *entity.Routes) error
	AddStudentToRoute(assignment *entity.RouteAssignment) error
	DeleteStudentFromRoute(routeNameUUID, studentUUID, schoolUUID string) error

	DeleteRoute(tx *sql.Tx, routenameUUID, schoolUUID string) error 
	DeleteRouteAssignments(tx *sql.Tx, routenameUUID, schoolUUID string) error
	RouteExists(tx *sql.Tx, routenameUUID, schoolUUID string) (bool, error)

	BeginTransaction() (*sql.Tx, error)

	IsDriverAssigned(tx *sql.Tx, driverUUID string) (bool, error)
	IsStudentAssigned(tx *sql.Tx, studentUUID string) (bool, error)
	GetDriverUUIDByRouteName(routeNameUUID string) (string, error)
	ValidateDriverVehicle(driverUUID string) (bool, error)
	CountAssignedStudentsByDriver(tx *sql.Tx, driverUUID string) (int, error)
	GetVehicleSeatsByDriver(tx *sql.Tx, driverUUID string) (int, error)
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
			ra.route_assignment_uuid,
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
			&route.RouteAssignmentUUID,
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
			r.route_assignment_uuid,
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
            route_assignment_uuid,
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
		assignment.RouteAssignmentUUID,
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

func (r *routeRepository) UpdateRouteDetails(route *entity.Routes) error {
    // Pastikan parameter yang dikirim sudah benar
    _, err := r.DB.Exec(`
        UPDATE routes 
        SET route_name = $1, route_description = $2, updated_by = $3, updated_at = $4
        WHERE route_name_uuid = $5 AND school_uuid = $6
    `, route.RouteName, route.RouteDescription, route.UpdatedBy, route.UpdatedAt, route.RouteNameUUID, route.SchoolUUID)
    
    if err != nil {
        log.Printf("Error updating route details: %v", err)
    }
    return err
}

func (r *routeRepository) AddStudentToRoute(assignment *entity.RouteAssignment) error {
    _, err := r.DB.Exec(`
        INSERT INTO route_assignment (
            route_id,
            route_assignment_uuid,
            driver_uuid,
            student_uuid,
            student_order,
            school_uuid,
            route_name_uuid,
            created_at,
            created_by
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `,
        assignment.RouteID,               // $1
        assignment.RouteAssignmentUUID,   // $2
        assignment.DriverUUID,            // $3 (pastikan diisi dengan benar)
        assignment.StudentUUID,           // $4
        assignment.StudentOrder,          // $5
        assignment.SchoolUUID,            // $6
        assignment.RouteNameUUID,         // $7
        assignment.CreatedAt,             // $8
        assignment.CreatedBy,             // $9
    )
    if err != nil {
        log.Printf("Error adding student to route: %v", err)
    }
    return err
}

func (r *routeRepository) DeleteStudentFromRoute(routeNameUUID, studentUUID, schoolUUID string) error {
    _, err := r.DB.Exec(`
        DELETE FROM route_assignment 
        WHERE route_name_uuid = $1 AND student_uuid = $2 AND school_uuid = $3
    `, routeNameUUID, studentUUID, schoolUUID)
    if err != nil {
        log.Printf("Error deleting student from route: %v", err)
    }
    return err
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

func (r *routeRepository) GetRouteAndStudentUUID(routeNameUUID string) (string, error) {
	var driverUUID *string
	query := `
		SELECT 
			ra.route_assignment_uuid,
			ra.student_uuid
		FROM route_assignment ra
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

func (r *routeRepository) GetVehicleSeatsByDriver(tx *sql.Tx, driverUUID string) (int, error) {
    var vehicleSeats int
    query := `
        SELECT v.vehicle_seats
        FROM vehicles v
        WHERE v.driver_uuid = $1
    `
    err := tx.QueryRow(query, driverUUID).Scan(&vehicleSeats)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, fmt.Errorf("vehicle with driver UUID %s not found", driverUUID)
        }
        return 0, fmt.Errorf("error fetching vehicle seats: %w", err)
    }
    return vehicleSeats, nil
}

func (r *routeRepository) CountAssignedStudentsByDriver(tx *sql.Tx, driverUUID string) (int, error) {
    var count int
    query := `
        SELECT COUNT(*)
        FROM route_assignment
        WHERE driver_uuid = $1
    `
    err := tx.QueryRow(query, driverUUID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("error counting assigned students: %w", err)
    }
    return count, nil
}
