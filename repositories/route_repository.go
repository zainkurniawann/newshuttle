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
	"github.com/umahmood/haversine"
)

type RouteRepositoryInterface interface {
	CountRoutesBySchool(schoolUUID string) (int, error)
	CalculateTotalDistance(driverStart [2]float64, students [][2]float64, school [2]float64) float64

	FetchAllRoutesByAS(offset, limit int, sortField, sortDirection, schoolUUID string) ([]dto.RoutesResponseDTO, error)
	FetchAllRouteAssignments(page, limit int) ([]dto.RoutesResponseDTO, int, error)
	FetchSpecRouteByAS(routeNameUUID, driverUUID string) ([]entity.RouteAssignment, error)
	FetchAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)

	AddRoutes(tx *sql.Tx, route entity.Routes) (string, error)
	AddRouteAssignment(tx *sql.Tx, assignment entity.RouteAssignment) error

	UpdateRouteDetails(route *entity.Routes) error
	AddStudentToRoute(assignment *entity.RouteAssignment) error
	UpdateStudentOrder(routeNameUUID string, assignment *entity.RouteAssignment, studentUUID string) error
	UpdateStudentOrderByDriver(studentUUID string, newOrder int) error
	GetMaxStudentOrder(routeNameUUID, schoolUUID string) (int, error)
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

func (r *routeRepository) CalculateTotalDistance(driverStart [2]float64, students [][2]float64, school [2]float64) float64 {
	totalDistance := 0.0
	currentLocation := driverStart

	// Hitung jarak dari driver ke siswa 1, siswa 1 ke siswa 2, dan seterusnya
	for _, student := range students {
		start := haversine.Coord{Lat: currentLocation[0], Lon: currentLocation[1]}
		end := haversine.Coord{Lat: student[0], Lon: student[1]}
		distance, _ := haversine.Distance(start, end)

		totalDistance += distance
		currentLocation = student // Update lokasi saat ini ke lokasi siswa yang baru
	}

	// Hitung jarak dari siswa terakhir ke sekolah
	start := haversine.Coord{Lat: currentLocation[0], Lon: currentLocation[1]}
	end := haversine.Coord{Lat: school[0], Lon: school[1]}
	distance, _ := haversine.Distance(start, end)
	totalDistance += distance

	// Mengonversi meter ke kilometer
	totalDistanceInKm := totalDistance / 1000 // hasil dalam kilometer

	// Print atau debug hasil sementara untuk memastikan
	fmt.Printf("Total distance (in meters): %f\n", totalDistance)
	fmt.Printf("Total distance (in kilometers): %f\n", totalDistanceInKm)

	// Mengembalikan hasil dengan dua angka desimal
	return totalDistanceInKm
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

func (r *routeRepository) FetchAllRouteAssignments(page, limit int) ([]dto.RoutesResponseDTO, int, error) {
    log.Println("Starting FetchAllRouteAssignments repository function")

    offset := (page - 1) * limit
    log.Printf("Pagination parameters - Page: %d, Limit: %d, Offset: %d\n", page, limit, offset)

    query := `
        SELECT 
            r.route_name_uuid,
            r.route_name,
            r.route_description,
            ra.route_assignment_uuid,
            ra.driver_uuid,
            dd.user_first_name AS driver_first_name,
            dd.user_last_name AS driver_last_name,
            ra.student_uuid,
            s.student_first_name,
            s.student_last_name,
            s.student_status,
			s.student_pickup_point,
            ra.student_order
        FROM route_assignment ra
        LEFT JOIN routes r ON ra.route_name_uuid = r.route_name_uuid
        LEFT JOIN driver_details dd ON ra.driver_uuid = dd.user_uuid
        LEFT JOIN students s ON ra.student_uuid = s.student_uuid
		LIMIT $1 OFFSET $2
    `

    log.Println("Executing query to fetch route assignments")
    rows, err := r.DB.Query(query, limit, offset)
    if err != nil {
        log.Printf("Error executing query: %v\n", err)
        return nil, 0, fmt.Errorf("failed to fetch route assignments: %w", err)
    }
    defer rows.Close()
    log.Println("Query executed successfully")

    // Struktur data mentah
    routeAssignments := make(map[string]*dto.RoutesResponseDTO)

    log.Println("Starting to process query result rows")
    for rows.Next() {
        var (
            routeNameUUID, routeName, routeDescription, routeAssignmentUUID string
            driverUUID, driverFirstName, driverLastName                     string
            studentUUID, studentFirstName, studentLastName                  string
            studentStatus, studentPickupPoint                               string
            studentOrder                                                    int
        )

        // Scan data hasil query
        err := rows.Scan(
			&routeNameUUID,
			&routeName,
			&routeDescription,
			&routeAssignmentUUID,
			&driverUUID,
			&driverFirstName,
			&driverLastName,
			&studentUUID,
			&studentFirstName,
			&studentLastName,
			&studentStatus,
			&studentPickupPoint,
			&studentOrder,
		)		
        if err != nil {
            log.Printf("Error scanning row: %v\n", err)
            return nil, 0, fmt.Errorf("failed to scan data: %w", err)
        }

        // Masukkan data ke dalam struktur
        if _, exists := routeAssignments[routeNameUUID]; !exists {
            log.Printf("Adding new route: %s\n", routeNameUUID)
            routeAssignments[routeNameUUID] = &dto.RoutesResponseDTO{
                RouteNameUUID:   routeNameUUID,
                RouteName:       routeName,
                RouteDescription: routeDescription,
                RouteAssignment:  []dto.RouteAssignmentResponseDTO{},
            }
        }

        currentRoute := routeAssignments[routeNameUUID]

        found := false
        for i := range currentRoute.RouteAssignment {
            if currentRoute.RouteAssignment[i].DriverUUID == driverUUID {
                log.Printf("Adding student to existing driver: %s\n", driverUUID)
                currentRoute.RouteAssignment[i].Students = append(currentRoute.RouteAssignment[i].Students, dto.StudentDTO{
                    StudentUUID:         studentUUID,
                    StudentFirstName:    studentFirstName,
                    StudentLastName:     studentLastName,
                    StudentStatus:       studentStatus,
					StudentPickupPoint:  studentPickupPoint,
                    StudentOrder:        studentOrder,
                    RouteAssignmentUUID: routeAssignmentUUID,
                })
                found = true
                break
            }
        }

        if !found {
            log.Printf("Adding new driver: %s\n", driverUUID)
            currentRoute.RouteAssignment = append(currentRoute.RouteAssignment, dto.RouteAssignmentResponseDTO{
                DriverUUID:      driverUUID,
                DriverFirstName: driverFirstName,
                DriverLastName:  driverLastName,
                Students: []dto.StudentDTO{
                    {
                        StudentUUID:         studentUUID,
                        StudentFirstName:    studentFirstName,
                        StudentLastName:     studentLastName,
                        StudentStatus:       studentStatus,
						StudentPickupPoint:  studentPickupPoint,
                        StudentOrder:        studentOrder,
                        RouteAssignmentUUID: routeAssignmentUUID,
                    },
                },
            })
			log.Printf("Scanned data: studentPickupPoint=%s\n", studentPickupPoint)
        }
    }

    // Hitung total items
    log.Println("Calculating total items")
    var totalItems int
    countQuery := "SELECT COUNT(*) FROM route_assignment"
    if err := r.DB.QueryRow(countQuery).Scan(&totalItems); err != nil {
        log.Printf("Error fetching total count: %v\n", err)
        return nil, 0, fmt.Errorf("failed to fetch total count: %w", err)
    }
    log.Printf("Total items: %d\n", totalItems)

    // Ubah map menjadi slice
    log.Println("Converting route assignments map to slice")
    result := make([]dto.RoutesResponseDTO, 0, len(routeAssignments))
    for _, value := range routeAssignments {
        result = append(result, *value)
    }

    log.Println("FetchAllRouteAssignments completed successfully")
    return result, totalItems, nil
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
	log.Println("Fetching routes for driver:", driverUUID)
	query := `
		SELECT
			r.route_assignment_uuid,
			r.student_uuid,
			r.driver_uuid,
			r.school_uuid,
			r.student_order,
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
		log.Println("Error fetching routes:", err)
		return nil, err
	}
	log.Println("Routes fetched successfully:", len(routes))
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

func (r *routeRepository) UpdateStudentOrder(routeNameUUID string, assignment *entity.RouteAssignment, studentUUID string) error {
    log.Println("Updating student order for RouteNameUUID:", assignment.RouteNameUUID)
    log.Printf("New student order: %d, routeNameUUID: %s, studentUUID: %s\n", assignment.StudentOrder, routeNameUUID, studentUUID)

    query := `
        UPDATE route_assignment
        SET student_order = $1
        WHERE route_name_uuid = $2 AND student_uuid = $3
    `

    _, err := r.DB.Exec(query, assignment.StudentOrder, routeNameUUID, studentUUID)
    if err != nil {
        log.Printf("Error executing query: %v\n", err)
        return fmt.Errorf("Gagal memperbarui student_order untuk student %s: %w", studentUUID, err)
    }

    log.Println("Student order successfully updated for RouteNameUUID:", routeNameUUID)
    return nil
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

func (r *routeRepository) GetMaxStudentOrder(routeNameUUID, schoolUUID string) (int, error) {
	query := `SELECT COALESCE(MAX(student_order), 0) FROM route_assignment WHERE route_name_uuid = $1 AND school_uuid = $2`
	var maxOrder int
	err := r.DB.QueryRow(query, routeNameUUID, schoolUUID).Scan(&maxOrder)
	if err != nil {
		return 0, err
	}
	return maxOrder, nil
}

func (repo *routeRepository) UpdateStudentOrderByDriver(studentUUID string, newOrder int) error {
	tx, err := repo.DB.Beginx()
	if err != nil {
		log.Println("Error starting transaction:", err)
		return err
	}

	// Validasi apakah student_uuid ada di database
	var exists bool
	err = tx.Get(&exists, "SELECT EXISTS(SELECT 1 FROM route_assignment WHERE student_uuid = $1)", studentUUID)
	if err != nil {
		tx.Rollback()
		log.Println("Error checking student_uuid existence:", err)
		return err
	}
	if !exists {
		tx.Rollback()
		log.Println("student_uuid not found:", studentUUID)
		return fmt.Errorf("student_uuid not found")
	}

	// Get the current order of the student
	var currentOrder int
	err = tx.Get(&currentOrder, "SELECT student_order FROM route_assignment WHERE student_uuid = $1", studentUUID)
	if err != nil {
		tx.Rollback()
		log.Println("Error getting current student order:", err)
		return err
	}

	// Shift the orders of other students
	if newOrder < currentOrder {
		_, err = tx.Exec("UPDATE route_assignment SET student_order = student_order + 1 WHERE student_order >= $1 AND student_order < $2", newOrder, currentOrder)
	} else {
		_, err = tx.Exec("UPDATE route_assignment SET student_order = student_order - 1 WHERE student_order > $1 AND student_order <= $2", currentOrder, newOrder)
	}
	if err != nil {
		tx.Rollback()
		log.Println("Error shifting student orders:", err)
		return err
	}

	// Update the student's order
	_, err = tx.Exec("UPDATE route_assignment SET student_order = $1 WHERE student_uuid = $2", newOrder, studentUUID)
	if err != nil {
		tx.Rollback()
		log.Println("Error updating student order:", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		return err
	}

	return nil
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
