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
	GetSpecVehicle(uuid string) (dto.VehicleResponseDTO, error)
	GetSpecVehicleForPermittedSchool(id string) (dto.VehicleResponseDTO, error)
	GetAllVehicles(page, limit int, sortField, sortDirection string) ([]dto.VehicleResponseDTO, int, error)
	GetAllVehiclesForPermittedSchool(page, limit int, sortField, sortDirection string) ([]dto.VehicleResponseDTO, int, error)
	AddVehicle(req dto.VehicleRequestDTO) error
	AddSchoolVehicleWithDriver(vehicle dto.VehicleDriverRequestDTO, driver dto.DriverDetailsRequestsDTO, schoolUUID string, username string) error
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

func (service *VehicleService) AddSchoolVehicleWithDriver(vehicle dto.VehicleDriverRequestDTO, driver dto.DriverDetailsRequestsDTO, schoolUUID string, username string) error {
	var driverID uuid.UUID

	// Periksa apakah email driver sudah ada di database
	driverExists, err := service.userRepository.CheckEmailExist("", vehicle.Driver.Email)
	if err != nil {
		return err
	}

	// Mulai transaksi
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return err
	}

	var transactionError error
	defer func() {
		if transactionError != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if !driverExists {
		// Jika driver belum ada, tambahkan data driver baru
		newDriver := &dto.UserRequestsDTO{
			Username:  vehicle.Driver.Username,
			FirstName: vehicle.Driver.FirstName,
			LastName:  vehicle.Driver.LastName,
			Gender:    vehicle.Driver.Gender,
			Email:     vehicle.Driver.Email,
			Password:  vehicle.Driver.Password,
			Role:      dto.Role(entity.Driver),
			RoleCode:  "D",
			Phone:     vehicle.Driver.Phone,
			Address:   vehicle.Driver.Address,
		}

		driverID, err = service.userService.AddUser(*newDriver, username)
		if err != nil {
			transactionError = err
			return transactionError
		}
	} else {
		// Jika driver sudah ada, ambil UUID-nya
		driverID, err = service.userRepository.FetchUUIDByEmail(vehicle.Driver.Email)
		if err != nil {
			transactionError = err
			return transactionError
		}
	}

	// Membuat entitas kendaraan
	newVehicle := &entity.Vehicle{
		ID:            time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
        UUID:          uuid.New(),
        VehicleName:   vehicle.Vehicle.Name,
        VehicleNumber: vehicle.Vehicle.Number,
        VehicleType:   vehicle.Vehicle.Type,
        VehicleColor:  vehicle.Vehicle.Color,
        VehicleSeats:  vehicle.Vehicle.Seats,
        VehicleStatus: vehicle.Vehicle.Status,
    }

	// Simpan data kendaraan
	err = service.vehicleRepository.SaveSchoolVehicleWithDriver(tx, *newVehicle)
	if err != nil {
		transactionError = err
		return transactionError
	}

	schoolUUIDParsed, err := uuid.Parse(schoolUUID)
	if err != nil {
		return fmt.Errorf("invalid school UUID: %v", err)
	}

	// // Parse schoolUUID dan buat pointer ke uuid.UUID
	// schoolUUIDParsed := uuid.Must(uuid.Parse(schoolUUID))

	// Membuat entitas DriverDetails
	driverDetails := &entity.DriverDetails{
		UserUUID:    driverID,
		SchoolUUID:  &schoolUUIDParsed, // Menggunakan pointer
		VehicleUUID: &newVehicle.UUID,
		LicenseNumber: driver.LicenseNumber,
	}

	// Simpan data DriverDetails dengan transaksi
	err = service.userRepository.SaveDriverDetails(tx, *driverDetails, driverID, nil)
	if err != nil {
		transactionError = err
		return transactionError
	}

	return nil
}

func (service *VehicleService) UpdateVehicle(id string, req dto.VehicleRequestDTO, username string) error {
    log.Println("Start updating vehicle with ID:", id)

    parsedUUID, err := uuid.Parse(id)
    if err != nil {
        log.Println("Error parsing vehicle UUID:", err)
        return err
    }

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
    log.Println("Vehicle entity to be updated:", vehicle)

    if req.School != "" {
        schoolUUID, err := uuid.Parse(req.School)
        if err != nil {
            log.Println("Error parsing school UUID:", err)
            return err
        }
        vehicle.SchoolUUID = &schoolUUID
    } else {
        vehicle.SchoolUUID = nil
    }
    log.Println("School UUID set to:", vehicle.SchoolUUID)

    // Cek apakah nomor kendaraan sudah ada
    isExistingVehicleNumber, err := service.vehicleRepository.CheckVehicleNumberExists(id, vehicle.VehicleNumber)
    if err != nil {
        log.Println("Error checking if vehicle number exists:", err)
        return err
    }

    if isExistingVehicleNumber {
        log.Println("Vehicle number already exists")
        return errors.New("Vehicle number already exists", 400)
    }
    log.Println("Vehicle number is unique")

    // Update kendaraan
    err = service.vehicleRepository.UpdateVehicle(vehicle)
    if err != nil {
        log.Println("Error updating vehicle:", err)
        return err
    }

    log.Println("Vehicle updated successfully")
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