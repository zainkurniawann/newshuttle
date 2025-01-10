package repositories

import (
	"fmt"
	"log"
	"shuttle/models/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type VehicleRepositoryInterface interface {
	CountVehicles() (int, error)
	CheckVehicleNumberExists(uuid ,vehicleNumber string) (bool, error)

	FetchAllVehicles(offset, limit int, sortField, sortDirection string) ([]entity.Vehicle, map[string]entity.School, map[string]entity.DriverDetails, error)
	FetchAllVehiclesForPermittedSchool(offset, limit int, sortField, sortDirection, schoolUUID string) ([]entity.Vehicle, map[string]entity.School, map[string]entity.DriverDetails, error)
	CountVehiclesForPermittedSchool(schoolUUID string) (int, error)
	FetchSpecVehicle(uuid string) (entity.Vehicle, entity.School, entity.DriverDetails, error)
	FetchSpecVehicleForPermittedSchool(uuid string) (entity.Vehicle, entity.School, entity.DriverDetails, error)

	SaveVehicle(vehicle entity.Vehicle) error
	SaveSchoolVehicleWithDriver(tx *sqlx.Tx, vehicle entity.Vehicle) error
	UpdateVehicle(vehicle entity.Vehicle) error
	DeleteVehicle(vehicle entity.Vehicle) error
}

type VehicleRepository struct {
	db *sqlx.DB
}

func NewVehicleRepository(db *sqlx.DB) VehicleRepositoryInterface {
	return &VehicleRepository{
		db: db,
	}
}

func (repository *VehicleRepository) CountVehicles() (int, error) {
	var count int

	query := `
		SELECT COUNT(vehicle_id)
		FROM vehicles
		WHERE deleted_at IS NULL
	`

	if err := repository.db.Get(&count, query); err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *VehicleRepository) CheckVehicleNumberExists(uuid, vehicleNumber string) (bool, error) {
	var count int

	query := `
		SELECT COUNT(vehicle_id)
		FROM vehicles
		WHERE vehicle_number = $1 AND deleted_at IS NULL
	`

	if uuid != "" {
		query += ` AND vehicle_uuid != $2`
		if err := repository.db.Get(&count, query, vehicleNumber, uuid); err != nil {
			return false, err
		}
	} else {
		if err := repository.db.Get(&count, query, vehicleNumber); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (repository *VehicleRepository) FetchAllVehicles(offset, limit int, sortField, sortDirection string) ([]entity.Vehicle, map[string]entity.School, map[string]entity.DriverDetails, error) {
    var vehicles []entity.Vehicle
    var schoolsMap = make(map[string]entity.School)
	var driversMap = make(map[string]entity.DriverDetails)

    query := fmt.Sprintf(`
        SELECT 
            v.vehicle_uuid, v.school_uuid, COALESCE(v.driver_uuid, NULL) AS driver_uuid,
			v.vehicle_name, v.vehicle_number, 
            v.vehicle_type, v.vehicle_color, v.vehicle_seats, v.vehicle_status, 
            v.created_at, 
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_uuid
				END, 
				NULL
			) AS school_uuid,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_name
				END,
				'N/A'
			) AS school_name,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_uuid
				END,
				NULL
			) AS driver_uuid,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_first_name
				END,
				'N/A'
			) AS driver_first_name,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_last_name
				END,
				'N/A'
			) AS driver_last_name
        FROM vehicles v
        LEFT JOIN schools s ON v.school_uuid = s.school_uuid
		LEFT JOIN driver_details d ON v.driver_uuid = d.user_uuid
		LEFT JOIN users u ON d.user_uuid = u.user_uuid
        WHERE v.deleted_at IS NULL
        ORDER BY %s %s
        LIMIT $1 OFFSET $2
    `, sortField, sortDirection)

    rows, err := repository.db.Queryx(query, limit, offset)
    if err != nil {
        return nil, nil, nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var vehicle entity.Vehicle
        var school entity.School
		var driver entity.DriverDetails

        err := rows.Scan(
            &vehicle.UUID, &vehicle.SchoolUUID, &vehicle.DriverUUID, &vehicle.VehicleName, &vehicle.VehicleNumber,
            &vehicle.VehicleType, &vehicle.VehicleColor, &vehicle.VehicleSeats, &vehicle.VehicleStatus,
            &vehicle.CreatedAt, &school.UUID, &school.Name, &driver.UserUUID, &driver.FirstName, &driver.LastName,
        )
        if err != nil {
            return nil, nil, nil, err
        }

        vehicles = append(vehicles, vehicle)
		if vehicle.SchoolUUID != nil && *vehicle.SchoolUUID != uuid.Nil {
			schoolsMap[vehicle.SchoolUUID.String()] = school
		}
		if vehicle.DriverUUID != nil && *vehicle.DriverUUID != uuid.Nil {
			driversMap[vehicle.DriverUUID.String()] = driver
		}
    }

    return vehicles, schoolsMap, driversMap, nil
}

func (repository *VehicleRepository) FetchAllVehiclesForPermittedSchool(offset, limit int, sortField, sortDirection, schoolUUID string) ([]entity.Vehicle, map[string]entity.School, map[string]entity.DriverDetails, error) {
    var vehicles []entity.Vehicle
    var schoolsMap = make(map[string]entity.School)
	var driversMap = make(map[string]entity.DriverDetails)

    // Mengubah query untuk menyertakan schoolUUID
    query := fmt.Sprintf(`
        SELECT 
            v.vehicle_uuid, v.school_uuid, COALESCE(v.driver_uuid, NULL) AS driver_uuid,
			v.vehicle_name, v.vehicle_number, 
            v.vehicle_type, v.vehicle_color, v.vehicle_seats, v.vehicle_status, 
            v.created_at, 
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_uuid
				END, 
				NULL
			) AS school_uuid,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_name
				END,
				'N/A'
			) AS school_name,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_uuid
				END,
				NULL
			) AS driver_uuid,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_first_name
				END,
				'N/A'
			) AS driver_first_name,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_last_name
				END,
				'N/A'
			) AS driver_last_name
        FROM vehicles v
        LEFT JOIN schools s ON v.school_uuid = s.school_uuid
		LEFT JOIN driver_details d ON v.driver_uuid = d.user_uuid
		LEFT JOIN users u ON d.user_uuid = u.user_uuid
        WHERE v.school_uuid = $1 AND v.deleted_at IS NULL
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortField, sortDirection)

    // Menggunakan schoolUUID sebagai parameter pertama dalam query
    rows, err := repository.db.Queryx(query, schoolUUID, limit, offset)
    if err != nil {
        return nil, nil, nil, err
    }
    defer rows.Close()

    // Proses hasil query
    for rows.Next() {
        var vehicle entity.Vehicle
        var school entity.School
		var driver entity.DriverDetails

        err := rows.Scan(
            &vehicle.UUID, &vehicle.SchoolUUID, &vehicle.DriverUUID, &vehicle.VehicleName, &vehicle.VehicleNumber,
            &vehicle.VehicleType, &vehicle.VehicleColor, &vehicle.VehicleSeats, &vehicle.VehicleStatus,
            &vehicle.CreatedAt, &school.UUID, &school.Name, &driver.UserUUID, &driver.FirstName, &driver.LastName,
        )
        if err != nil {
            return nil, nil, nil, err
        }

        vehicles = append(vehicles, vehicle)
		if vehicle.SchoolUUID != nil && *vehicle.SchoolUUID != uuid.Nil {
			schoolsMap[vehicle.SchoolUUID.String()] = school
		}
		if vehicle.DriverUUID != nil && *vehicle.DriverUUID != uuid.Nil {
			driversMap[vehicle.DriverUUID.String()] = driver
		}
    }

    return vehicles, schoolsMap, driversMap, nil
}

func (repository *VehicleRepository) CountVehiclesForPermittedSchool(schoolUUID string) (int, error) {
	var count int

	// Menggunakan schoolUUID sebagai parameter untuk menghitung kendaraan berdasarkan sekolah tertentu
	query := `
		SELECT COUNT(vehicle_id)
		FROM vehicles
		WHERE school_uuid = $1 AND deleted_at IS NULL
	`

	// Mengambil data count berdasarkan query
	if err := repository.db.Get(&count, query, schoolUUID); err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *VehicleRepository) FetchSpecVehicle(uuid string) (entity.Vehicle, entity.School, entity.DriverDetails, error) {
	var vehicle entity.Vehicle
	var school entity.School
	var driver entity.DriverDetails

	query := `
		SELECT
			v.vehicle_uuid, v.school_uuid, v.driver_uuid, v.vehicle_name, v.vehicle_number,
			v.vehicle_type, v.vehicle_color, v.vehicle_seats, v.vehicle_status,
			v.created_at, v.created_by, v.updated_at, v.updated_by,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_uuid
				END, 
				NULL
			) AS school_uuid,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_name
				END,
				'N/A'
			) AS school_name,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_uuid
				END,
				NULL
			) AS driver_uuid,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_first_name
				END,
				'N/A'
			) AS driver_first_name,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_last_name
				END,
				'N/A'
			) AS driver_last_name
		FROM vehicles v
		LEFT JOIN schools s ON v.school_uuid = s.school_uuid
		LEFT JOIN driver_details d ON v.driver_uuid = d.user_uuid
		LEFT JOIN users u ON d.user_uuid = u.user_uuid
		WHERE v.deleted_at IS NULL AND v.vehicle_uuid = $1
	`

	err := repository.db.QueryRowx(query, uuid).Scan(
		&vehicle.UUID, &vehicle.SchoolUUID, &vehicle.DriverUUID, &vehicle.VehicleName, &vehicle.VehicleNumber,
		&vehicle.VehicleType, &vehicle.VehicleColor, &vehicle.VehicleSeats, &vehicle.VehicleStatus,
		&vehicle.CreatedAt, &vehicle.CreatedBy, &vehicle.UpdatedAt, &vehicle.UpdatedBy,
		&school.UUID, &school.Name, &driver.UserUUID, &driver.FirstName, &driver.LastName,
	)
	if err != nil {
		return entity.Vehicle{}, entity.School{}, entity.DriverDetails{}, err
	}

	return vehicle, school, driver, nil
}

func (repository *VehicleRepository) FetchSpecVehicleForPermittedSchool(uuid string) (entity.Vehicle, entity.School, entity.DriverDetails, error) {
	var vehicle entity.Vehicle
	var school entity.School
	var driver entity.DriverDetails

	query := `
		SELECT
			v.vehicle_uuid, v.school_uuid, v.driver_uuid, v.vehicle_name, v.vehicle_number,
			v.vehicle_type, v.vehicle_color, v.vehicle_seats, v.vehicle_status,
			v.created_at, v.created_by, v.updated_at, v.updated_by,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_uuid
				END, 
				NULL
			) AS school_uuid,
			COALESCE(
				CASE
					WHEN s.deleted_at IS NULL THEN s.school_name
				END,
				'N/A'
			) AS school_name,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_uuid
				END,
				NULL
			) AS driver_uuid,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_first_name
				END,
				'N/A'
			) AS driver_first_name,
			COALESCE(
				CASE
					WHEN u.deleted_at IS NULL THEN d.user_last_name
				END,
				'N/A'
			) AS driver_last_name
		FROM vehicles v
		LEFT JOIN schools s ON v.school_uuid = s.school_uuid
		LEFT JOIN driver_details d ON v.driver_uuid = d.user_uuid
		LEFT JOIN users u ON d.user_uuid = u.user_uuid
		WHERE v.deleted_at IS NULL AND v.vehicle_uuid = $1
	`

	err := repository.db.QueryRowx(query, uuid).Scan(
		&vehicle.UUID, &vehicle.SchoolUUID, &vehicle.DriverUUID, &vehicle.VehicleName, &vehicle.VehicleNumber,
		&vehicle.VehicleType, &vehicle.VehicleColor, &vehicle.VehicleSeats, &vehicle.VehicleStatus,
		&vehicle.CreatedAt, &vehicle.CreatedBy, &vehicle.UpdatedAt, &vehicle.UpdatedBy,
		&school.UUID, &school.Name, &driver.UserUUID, &driver.FirstName, &driver.LastName,
	)
	if err != nil {
		return entity.Vehicle{}, entity.School{}, entity.DriverDetails{}, err
	}

	return vehicle, school, driver, nil
}

func (repository *VehicleRepository) SaveVehicle(vehicle entity.Vehicle) error {
	query := `
		INSERT INTO vehicles (vehicle_id, vehicle_uuid, school_uuid, vehicle_name, vehicle_number, vehicle_type, vehicle_color, vehicle_seats, vehicle_status, created_by)
		VALUES (:vehicle_id, :vehicle_uuid, :school_uuid, :vehicle_name, :vehicle_number, :vehicle_type, :vehicle_color, :vehicle_seats, :vehicle_status, :created_by)
	`

	_, err := repository.db.NamedExec(query, vehicle)
	if err != nil {
		return err
	}

	return nil
}

func (repository *VehicleRepository) SaveSchoolVehicleWithDriver(tx *sqlx.Tx, vehicle entity.Vehicle) error {
    log.Println("Inserting vehicle into database:", vehicle)

    query := `
        INSERT INTO vehicles (vehicle_id, vehicle_uuid, school_uuid, vehicle_name, vehicle_number, vehicle_type, vehicle_color, vehicle_seats, vehicle_status, created_by)
        VALUES (:vehicle_id, :vehicle_uuid, :school_uuid, :vehicle_name, :vehicle_number, :vehicle_type, :vehicle_color, :vehicle_seats, :vehicle_status, :created_by)
    `
    log.Printf("SQL query to insert vehicle: %s\n", query)

    _, err := repository.db.NamedExec(query, vehicle)
    if err != nil {
        log.Println("Error inserting vehicle:", err)
        return err
    }
    log.Println("Vehicle inserted successfully into the database")
    return nil
}

func (repository *VehicleRepository) UpdateVehicle(vehicle entity.Vehicle) error {
	query := `
		UPDATE vehicles
		SET school_uuid = :school_uuid, vehicle_name = :vehicle_name, vehicle_number = :vehicle_number, vehicle_type = :vehicle_type, vehicle_color = :vehicle_color,
		vehicle_seats = :vehicle_seats, vehicle_status = :vehicle_status, updated_at = :updated_at, updated_by = :updated_by
		WHERE vehicle_uuid = :vehicle_uuid
	`

	_, err := repository.db.NamedExec(query, vehicle)
	if err != nil {
		return err
	}

	return nil
}

func (repository *VehicleRepository) DeleteVehicle(vehicle entity.Vehicle) error {
	query := `
		UPDATE vehicles
		SET deleted_at = :deleted_at, deleted_by = :deleted_by
		WHERE vehicle_uuid = :vehicle_uuid
	`

	_, err := repository.db.NamedExec(query, vehicle)
	if err != nil {
		return err
	}

	return nil
}