package services

import (
	"fmt"
	"log"
	"shuttle/errors"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
)

type VehicleServiceInterface interface {
	GetAvailableVehicles() ([]dto.VehicleResponseDTO, error)
	GetAvailableSchoolVehicles(schoolUUID string) ([]dto.VehicleResponseDTO, error)
	GetSpecVehicle(uuid string) (dto.VehicleResponseDTO, error)
	GetSpecVehicleForPermittedSchool(id string) (dto.VehicleResponseDTO, error)
	GetAllVehicles(page, limit int, sortField, sortDirection string) ([]dto.VehicleResponseDTO, int, error)
	GetAllVehiclesForPermittedSchool(page, limit int, sortField, sortDirection string) ([]dto.VehicleResponseDTO, int, error)
	AddVehicle(req dto.VehicleRequestDTO) error
	//AddSchoolVehicleWithDriver(vehicle dto.VehicleDriverRequestDTO, driver dto.DriverDetailsRequestsDTO, schoolUUID string, username string) error
	AddVehicleForPermittedSchool(req dto.VehicleRequestDTO, role, schoolUUID string) error
	UpdateVehicle(id string, req dto.VehicleRequestDTO, username string) error
	DeleteVehicle(id string, username string) error
}

type VehicleService struct {
	userService       UserServiceInterface
	vehicleRepository repositories.VehicleRepositoryInterface
	userRepository    repositories.UserRepositoryInterface
}

func NewVehicleService(vehicleRepository repositories.VehicleRepositoryInterface) VehicleService {
	return VehicleService{
		vehicleRepository: vehicleRepository,
	}
}

func (service *VehicleService) GetAllVehicles(page, limit int, sortField, sortDirection string) ([]dto.VehicleResponseDTO, int, error) {
	offset := (page - 1) * limit

	vehicles, school, driver, err := service.vehicleRepository.FetchAllVehicles(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.vehicleRepository.CountVehicles()
	if err != nil {
		return nil, 0, err
	}

	var vehiclesDTO []dto.VehicleResponseDTO
	for _, vehicle := range vehicles {

		var schoolName string
		if vehicle.SchoolUUID == nil || school[vehicle.SchoolUUID.String()].UUID == uuid.Nil {
			schoolName = "N/A"
		} else {
			schoolName = school[vehicle.SchoolUUID.String()].Name
		}

		var driverName string
		if vehicle.DriverUUID == nil || driver[vehicle.DriverUUID.String()].UserUUID == uuid.Nil {
			driverName = "N/A"
		} else {
			driverName = driver[vehicle.DriverUUID.String()].FirstName + " " + driver[vehicle.DriverUUID.String()].LastName
		}

		vehiclesDTO = append(vehiclesDTO, dto.VehicleResponseDTO{
			UUID:       vehicle.UUID.String(),
			SchoolName: schoolName,
			DriverName: driverName,
			Name:       vehicle.VehicleName,
			Number:     vehicle.VehicleNumber,
			Type:       vehicle.VehicleType,
			Color:      vehicle.VehicleColor,
			Seats:      vehicle.VehicleSeats,
			Status:     vehicle.VehicleStatus,
			CreatedAt:  safeTimeFormat(vehicle.CreatedAt),
		})
	}

	return vehiclesDTO, total, nil
}

func (service *VehicleService) GetAllVehiclesForPermittedSchool(page, limit int, sortField, sortDirection, schoolUUID string) ([]dto.VehicleResponseDTO, int, error) {
    offset := (page - 1) * limit

    // Modifikasi query untuk memasukkan schoolUUID
    vehicles, school, driver, err := service.vehicleRepository.FetchAllVehiclesForPermittedSchool(offset, limit, sortField, sortDirection, schoolUUID)
    if err != nil {
        return nil, 0, err
    }

	total, err := service.vehicleRepository.CountVehiclesForPermittedSchool(schoolUUID)
    if err != nil {
        return nil, 0, err
    }

    var vehiclesDTO []dto.VehicleResponseDTO
    for _, vehicle := range vehicles {

        var schoolName string
        if vehicle.SchoolUUID == nil || school[vehicle.SchoolUUID.String()].UUID == uuid.Nil {
            schoolName = "N/A"
        } else {
            schoolName = school[vehicle.SchoolUUID.String()].Name
        }

        var driverName string
        if vehicle.DriverUUID == nil || driver[vehicle.DriverUUID.String()].UserUUID == uuid.Nil {
            driverName = "N/A"
        } else {
            driverName = driver[vehicle.DriverUUID.String()].FirstName + " " + driver[vehicle.DriverUUID.String()].LastName
        }

        vehiclesDTO = append(vehiclesDTO, dto.VehicleResponseDTO{
            UUID:       vehicle.UUID.String(),
            SchoolName: schoolName,
            DriverName: driverName,
            Name:       vehicle.VehicleName,
            Number:     vehicle.VehicleNumber,
            Type:       vehicle.VehicleType,
            Color:      vehicle.VehicleColor,
            Seats:      vehicle.VehicleSeats,
            Status:     vehicle.VehicleStatus,
            CreatedAt:  safeTimeFormat(vehicle.CreatedAt),
        })
    }

    return vehiclesDTO, total, nil
}

func (service *VehicleService) GetSpecVehicle(id string) (dto.VehicleResponseDTO, error) {
	vehicle, school, driver, err := service.vehicleRepository.FetchSpecVehicle(id)
	if err != nil {
		return dto.VehicleResponseDTO{}, err
	}

	var schoolUUID, schoolName string
	if vehicle.SchoolUUID == nil {
		schoolUUID = "N/A"
		schoolName = "N/A"
	} else if vehicle.SchoolUUID != nil {
		schoolUUID = vehicle.SchoolUUID.String()
		schoolName = school.Name
	}

	var driverUUID, driverName string
	if driver.UserUUID == uuid.Nil {
		driverUUID = "N/A"
		driverName = "N/A"
	} else if driver.UserUUID != uuid.Nil {
		driverUUID = vehicle.DriverUUID.String()
		driverName = driver.FirstName + " " + driver.LastName
	}

	vehicleDTO := dto.VehicleResponseDTO{
		UUID:       vehicle.UUID.String(),
		SchoolUUID: schoolUUID,
		SchoolName: schoolName,
		DriverUUID: driverUUID,
		DriverName: driverName,
		Name:       vehicle.VehicleName,
		Number:     vehicle.VehicleNumber,
		Type:       vehicle.VehicleType,
		Color:      vehicle.VehicleColor,
		Seats:      vehicle.VehicleSeats,
		Status:     vehicle.VehicleStatus,
		CreatedAt:  safeTimeFormat(vehicle.CreatedAt),
		CreatedBy:  safeStringFormat(vehicle.CreatedBy),
		UpdatedAt:  safeTimeFormat(vehicle.UpdatedAt),
		UpdatedBy:  safeStringFormat(vehicle.UpdatedBy),
	}

	return vehicleDTO, nil
}

func (service *VehicleService) GetSpecVehicleForPermittedSchool(id string) (dto.VehicleResponseDTO, error) {
	vehicle, school, driver, err := service.vehicleRepository.FetchSpecVehicleForPermittedSchool(id)
	if err != nil {
		return dto.VehicleResponseDTO{}, err
	}

	var schoolUUID, schoolName string
	if vehicle.SchoolUUID == nil {
		schoolUUID = "N/A"
		schoolName = "N/A"
	} else if vehicle.SchoolUUID != nil {
		schoolUUID = vehicle.SchoolUUID.String()
		schoolName = school.Name
	}

	var driverUUID, driverName string
	if driver.UserUUID == uuid.Nil {
		driverUUID = "N/A"
		driverName = "N/A"
	} else if driver.UserUUID != uuid.Nil {
		driverUUID = vehicle.DriverUUID.String()
		driverName = driver.FirstName + " " + driver.LastName
	}

	vehicleDTO := dto.VehicleResponseDTO{
		UUID:       vehicle.UUID.String(),
		SchoolUUID: schoolUUID,
		SchoolName: schoolName,
		DriverUUID: driverUUID,
		DriverName: driverName,
		Name:       vehicle.VehicleName,
		Number:     vehicle.VehicleNumber,
		Type:       vehicle.VehicleType,
		Color:      vehicle.VehicleColor,
		Seats:      vehicle.VehicleSeats,
		Status:     vehicle.VehicleStatus,
		CreatedAt:  safeTimeFormat(vehicle.CreatedAt),
		CreatedBy:  safeStringFormat(vehicle.CreatedBy),
		UpdatedAt:  safeTimeFormat(vehicle.UpdatedAt),
		UpdatedBy:  safeStringFormat(vehicle.UpdatedBy),
	}

	return vehicleDTO, nil
}

func (service *VehicleService) GetAvailableVehicles() ([]dto.VehicleResponseDTO, error) {
	// Fetch available vehicles from repository
	vehicles, err := service.vehicleRepository.FetchAvailableVehicle()
	if err != nil {
		log.Printf("Error fetching available vehicles: %v", err)  // Log error
		return nil, fmt.Errorf("gagal mengambil data kendaraan yang tersedia: %w", err)
	}

	// Convert the vehicle data to DTO format
	var vehicleDTOs []dto.VehicleResponseDTO
	for _, vehicle := range vehicles {
		vehicleDTO := dto.VehicleResponseDTO{
			UUID:          vehicle.UUID.String(),
			Name:   vehicle.VehicleName,
			Number: vehicle.VehicleNumber,
			Color:  vehicle.VehicleColor,
		}
		vehicleDTOs = append(vehicleDTOs, vehicleDTO)
	}

	return vehicleDTOs, nil
}

func (service *VehicleService) GetAvailableSchoolVehicles(schoolUUID string) ([]dto.VehicleResponseDTO, error) {
    // Fetch available vehicles from repository
    vehicles, err := service.vehicleRepository.FetchAvailableSchoolVehicle(schoolUUID)
    if err != nil {
        log.Printf("Error fetching available vehicles: %v", err)
        return nil, fmt.Errorf("gagal mengambil data kendaraan yang tersedia: %w", err)
    }

    // Convert the vehicle data to DTO format
    var vehicleDTOs []dto.VehicleResponseDTO
    for _, vehicle := range vehicles {
        vehicleDTO := dto.VehicleResponseDTO{
            UUID:   vehicle.UUID.String(),
            Name:   vehicle.VehicleName,
            Number: vehicle.VehicleNumber,
            Color:  vehicle.VehicleColor,
        }
        vehicleDTOs = append(vehicleDTOs, vehicleDTO)
    }

    return vehicleDTOs, nil
}

func (service *VehicleService) AddVehicle(req dto.VehicleRequestDTO) error {
	vehicle := entity.Vehicle{
		ID:            time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		UUID:          uuid.New(),
		VehicleName:   req.Name,
		VehicleNumber: req.Number,
		VehicleType:   req.Type,
		VehicleColor:  req.Color,
		VehicleSeats:  req.Seats,
		VehicleStatus: req.Status,
	}

	if req.School != "" {
		schoolUUID, err := uuid.Parse(req.School)
		if err != nil {
			return err
		}
		vehicle.SchoolUUID = &schoolUUID
	} else {
		vehicle.SchoolUUID = nil
	}

	isExistingVehicleNumber, err := service.vehicleRepository.CheckVehicleNumberExists("", vehicle.VehicleNumber)
	if err != nil {
		return err
	}

	if isExistingVehicleNumber {
		return errors.New("Vehicle number already exists", 400)
	}

	err = service.vehicleRepository.SaveVehicle(vehicle)
	if err != nil {
		return err
	}

	return nil
}

// func (service *VehicleService) AddSchoolVehicleWithDriver(vehicle dto.VehicleDriverRequestDTO, driver dto.DriverDetailsRequestsDTO, schoolUUID string, username string) error {
// 	var driverID uuid.UUID

// 	// Periksa apakah email driver sudah ada di database
// 	driverExists, err := service.userRepository.CheckEmailExist("", vehicle.Driver.Email)
// 	if err != nil {
// 		return err
// 	}

// 	// Mulai transaksi
// 	tx, err := service.userRepository.BeginTransaction()
// 	if err != nil {
// 		return err
// 	}

// 	var transactionError error
// 	defer func() {
// 		if transactionError != nil {
// 			tx.Rollback()
// 		} else {
// 			tx.Commit()
// 		}
// 	}()

// 	if !driverExists {
// 		// Jika driver belum ada, tambahkan data driver baru
// 		newDriver := &dto.UserRequestsDTO{
// 			Username:  vehicle.Driver.Username,
// 			FirstName: vehicle.Driver.FirstName,
// 			LastName:  vehicle.Driver.LastName,
// 			Gender:    vehicle.Driver.Gender,
// 			Email:     vehicle.Driver.Email,
// 			Password:  vehicle.Driver.Password,
// 			Role:      dto.Role(entity.Driver),
// 			RoleCode:  "D",
// 			Phone:     vehicle.Driver.Phone,
// 			Address:   vehicle.Driver.Address,
// 		}

// 		driverID, err = service.userService.AddUser(*newDriver, username)
// 		if err != nil {
// 			transactionError = err
// 			return transactionError
// 		}
// 	} else {
// 		// Jika driver sudah ada, ambil UUID-nya
// 		driverID, err = service.userRepository.FetchUUIDByEmail(vehicle.Driver.Email)
// 		if err != nil {
// 			transactionError = err
// 			return transactionError
// 		}
// 	}

// 	// Membuat entitas kendaraan
// 	newVehicle := &entity.Vehicle{
// 		ID:            time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
//         UUID:          uuid.New(),
//         VehicleName:   vehicle.Vehicle.Name,
//         VehicleNumber: vehicle.Vehicle.Number,
//         VehicleType:   vehicle.Vehicle.Type,
//         VehicleColor:  vehicle.Vehicle.Color,
//         VehicleSeats:  vehicle.Vehicle.Seats,
//         VehicleStatus: vehicle.Vehicle.Status,
//     }

// 	// Simpan data kendaraan
// 	err = service.vehicleRepository.SaveSchoolVehicleWithDriver(tx, *newVehicle)
// 	if err != nil {
// 		transactionError = err
// 		return transactionError
// 	}

// 	schoolUUIDParsed, err := uuid.Parse(schoolUUID)
// 	if err != nil {
// 		return fmt.Errorf("invalid school UUID: %v", err)
// 	}

// 	// // Parse schoolUUID dan buat pointer ke uuid.UUID
// 	// schoolUUIDParsed := uuid.Must(uuid.Parse(schoolUUID))

// 	// Membuat entitas DriverDetails
// 	driverDetails := &entity.DriverDetails{
// 		UserUUID:    driverID,
// 		SchoolUUID:  &schoolUUIDParsed, // Menggunakan pointer
// 		VehicleUUID: &newVehicle.UUID,
// 		LicenseNumber: driver.LicenseNumber,
// 	}

// 	// Simpan data DriverDetails dengan transaksi
// 	err = service.userRepository.SaveDriverDetails(tx, *driverDetails, driverID, nil)
// 	if err != nil {
// 		transactionError = err
// 		return transactionError
// 	}

// 	return nil
// }

func (service *VehicleService) AddVehicleForPermittedSchool(req dto.VehicleRequestDTO, role, schoolUUID string) error {
    log.Println("Start adding vehicle")

    // Membuat entitas vehicle
    vehicle := entity.Vehicle{
        ID:            time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
        UUID:          uuid.New(),
        VehicleName:   req.Name,
        VehicleNumber: req.Number,
        VehicleType:   req.Type,
        VehicleColor:  req.Color,
        VehicleSeats:  req.Seats,
        VehicleStatus: req.Status,
    }
    log.Printf("Vehicle entity created: %+v\n", vehicle)
    // Gunakan schoolUUID yang sudah ada di context
    if role == "AS" {
        log.Println("Role is schooladmin, using school_uuid from token")
        if schoolUUID != "" {
            schoolUUIDParsed, err := uuid.Parse(schoolUUID)
            if err != nil {
                log.Println("Error parsing school UUID:", err)
                return errors.New("Invalid school UUID", 400)
            }
            vehicle.SchoolUUID = &schoolUUIDParsed
            log.Println("School UUID parsed and assigned to vehicle:", schoolUUIDParsed)
        } else {
            log.Println("School UUID is required for schooladmin")
            return errors.New("School UUID is required for schooladmin", 400)
        }
    } else {
        log.Println("Non-schooladmin role, using provided school UUID or nil")
    }

    // Cek apakah vehicle number sudah ada
    log.Println("Checking if vehicle number exists:", vehicle.VehicleNumber)
    isExistingVehicleNumber, err := service.vehicleRepository.CheckVehicleNumberExists("", vehicle.VehicleNumber)
    if err != nil {
        log.Println("Error checking vehicle number:", err)
        return err
    }

    if isExistingVehicleNumber {
        log.Println("Vehicle number already exists:", vehicle.VehicleNumber)
        return errors.New("Vehicle number already exists", 400)
    }
    log.Println("Vehicle number is unique, proceeding to save")

    // Simpan kendaraan
    log.Println("Saving vehicle to database")
    err = service.vehicleRepository.SaveVehicleForPermittedSchool(vehicle)
    if err != nil {
        log.Println("Error saving vehicle:", err)
        return err
    }
    log.Println("Vehicle saved successfully")
    return nil
}

func (service *VehicleService) UpdateVehicle(id string, req dto.VehicleRequestDTO, username string) error {
    log.Println("Start updating vehicle with ID:", id)

    // Parsing UUID
    parsedUUID, err := uuid.Parse(id)
    if err != nil {
        log.Println("Error parsing vehicle UUID:", err)
        return err
    }
    log.Println("Parsed vehicle UUID successfully:", parsedUUID)

    // Mapping request data to entity
    vehicle := entity.Vehicle{
        UUID:          parsedUUID,
        VehicleName:   req.Name,
        VehicleNumber: req.Number,
        VehicleType:   req.Type,
        VehicleColor:  req.Color,
        VehicleSeats:  req.Seats,
        VehicleStatus: req.Status,
        UpdatedAt:     toNullTime(time.Now()),
        UpdatedBy:     toNullString(username),
    }
    log.Println("Vehicle entity constructed:", vehicle)

    // Parsing and setting School UUID if present
    if req.School != "" {
        log.Println("Parsing School UUID:", req.School)
        schoolUUID, err := uuid.Parse(req.School)
        if err != nil {
            log.Println("Error parsing school UUID:", err)
            return err
        }
        vehicle.SchoolUUID = &schoolUUID
        log.Println("School UUID parsed and set successfully:", schoolUUID)
    } else {
        vehicle.SchoolUUID = nil
        log.Println("School UUID not provided, set to nil")
    }

    // Check if vehicle number already exists
    log.Println("Checking if vehicle number exists for ID:", id, "and number:", vehicle.VehicleNumber)
    isExistingVehicleNumber, err := service.vehicleRepository.CheckVehicleNumberExists(id, vehicle.VehicleNumber)
    if err != nil {
        log.Println("Error checking vehicle number existence:", err)
        return err
    }

    if isExistingVehicleNumber {
        log.Println("Vehicle number already exists:", vehicle.VehicleNumber)
        return errors.New("Vehicle number already exists", 400)
    }
    log.Println("Vehicle number is unique:", vehicle.VehicleNumber)

    // Updating vehicle
    log.Println("Updating vehicle in repository with data:", vehicle)
    err = service.vehicleRepository.UpdateVehicle(vehicle)
    if err != nil {
        log.Println("Error updating vehicle:", err)
        return err
    }

    log.Println("Vehicle updated successfully with ID:", id)
    return nil
}

func (service *VehicleService) DeleteVehicle(id string, username string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	vehicle := entity.Vehicle{
		UUID:      parsedUUID,
		DeletedAt: toNullTime(time.Now()),
		DeletedBy: toNullString(username),
	}

	err = service.vehicleRepository.DeleteVehicle(vehicle)
	if err != nil {
		return err
	}

	return nil
}