package services

import (
	"database/sql"
	"fmt"
	"log"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
)

type ShuttleServiceInterface interface {
	GetShuttleTrackByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error)
	GetAllShuttleByParent(parentUUID uuid.UUID, page, limit int, sortField, sortDirection string) ([]dto.ShuttleAllResponse, int, error)
	GetAllShuttleByDriver(driverUUID uuid.UUID) ([]dto.ShuttleAllResponse, error)
	GetSpecShuttle(shuttleUUID uuid.UUID) ([]dto.ShuttleSpecResponse, error)
	AddShuttle(req dto.ShuttleRequest, driverUUID, createdBy string) error
	EditShuttleStatus(shuttleUUID, status string) error
}

type ShuttleService struct {
	shuttleRepository repositories.ShuttleRepositoryInterface
}

func NewShuttleService(shuttleRepository repositories.ShuttleRepositoryInterface) ShuttleServiceInterface {
	return &ShuttleService{
		shuttleRepository: shuttleRepository,
	}
}

func (s *ShuttleService) GetShuttleTrackByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error) {
	log.Println("Fetching shuttle track from repository for parentUUID:", parentUUID)

	// Panggil repository dengan parameter baru
	shuttles, err := s.shuttleRepository.FetchShuttleTrackByParent(parentUUID)
	if err != nil {
		log.Println("Error fetching shuttle track from repository:", err)
		return nil, err
	}

	log.Println("Fetched shuttle data:", shuttles)
	responses := make([]dto.ShuttleResponse, 0, len(shuttles))
	for _, shuttle := range shuttles {
		response := &dto.ShuttleResponse{
			StudentUUID:    shuttle.StudentUUID,
			ShuttleUUID:    shuttle.ShuttleUUID,
			StudentFirstName:    shuttle.StudentFirstName,
			StudentLastName: shuttle.StudentLastName,
			ParentUUID:     shuttle.ParentUUID,
			SchoolUUID:     shuttle.SchoolUUID,
			SchoolName:     shuttle.SchoolName,
			ShuttleStatus:  shuttle.ShuttleStatus,
			CreatedAt:      shuttle.CreatedAt,
			CurrentDate:    shuttle.CurrentDate,
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

func (s *ShuttleService) GetAllShuttleByParent(parentUUID uuid.UUID, page, limit int, sortField, sortDirection string) ([]dto.ShuttleAllResponse, int, error) {
    offset := (page - 1) * limit

    // Panggil repository untuk mendapatkan data dan total items
    shuttles, err := s.shuttleRepository.FetchAllShuttleByParent(offset, limit, sortField, sortDirection, parentUUID)
    if err != nil {
        return nil, 0, err
    }

    totalItems, err := s.shuttleRepository.CountShuttlesByParent(parentUUID)
    if err != nil {
        return nil, 0, err
    }

    return shuttles, totalItems, nil
}

func (s *ShuttleService) GetAllShuttleByDriver(driverUUID uuid.UUID) ([]dto.ShuttleAllResponse, error) {
	// Fetch data from the repository
	shuttles, err := s.shuttleRepository.FetchAllShuttleByDriver(driverUUID)
	if err != nil {
		return nil, err
	}

	// Transform the data if needed (DTO is already in the required format)
	responses := make([]dto.ShuttleAllResponse, len(shuttles))
	for i, shuttle := range shuttles {
		responses[i] = dto.ShuttleAllResponse{
			ShuttleUUID:    shuttle.ShuttleUUID,
			StudentUUID:     shuttle.StudentUUID,
			Status:          shuttle.Status,
			StudentFirstName: shuttle.StudentFirstName,
			StudentLastName:  shuttle.StudentLastName,
			StudentGrade:     shuttle.StudentGrade,
			StudentGender:    shuttle.StudentGender,
			ParentUUID:       shuttle.ParentUUID,
			SchoolUUID:       shuttle.SchoolUUID,
			SchoolName:       shuttle.SchoolName,
			CreatedAt:        shuttle.CreatedAt,
			UpdatedAt:        shuttle.UpdatedAt,
		}
	}

	return responses, nil
}

func (s *ShuttleService) GetSpecShuttle(shuttleUUID uuid.UUID) ([]dto.ShuttleSpecResponse, error) {
	log.Println("Fetching shuttle spec data from repository for shuttleUUID:", shuttleUUID)

	shuttles, err := s.shuttleRepository.GetSpecShuttle(shuttleUUID)
	if err != nil {
		log.Println("Error fetching shuttle data from repository:", err)
		return nil, fmt.Errorf("failed to fetch shuttle data: %w", err)
	}

	log.Println("Fetched shuttle data from repository:", shuttles)
	responses := make([]dto.ShuttleSpecResponse, 0, len(shuttles))
	for _, shuttle := range shuttles {
		log.Println("Processing shuttle:", shuttle)
		response := dto.ShuttleSpecResponse{
			StudentUUID:       shuttle.StudentUUID,
			ShuttleUUID:       shuttle.ShuttleUUID,
			StudentFirstName:  shuttle.StudentFirstName,
			StudentLastName:   shuttle.StudentLastName,
			StudentPickupPoint: shuttle.StudentPickupPoint,
			ParentUUID:        shuttle.ParentUUID,
			SchoolUUID:        shuttle.SchoolUUID,
			SchoolName:        shuttle.SchoolName,
			SchoolPoint:       shuttle.SchoolPoint,
			ShuttleStatus:     shuttle.ShuttleStatus,
			CreatedAt:         shuttle.CreatedAt,
			CurrentDate:       shuttle.CurrentDate,
			DriverUUID:        shuttle.DriverUUID,
			DriverUsername:    shuttle.DriverUsername,
			DriverFirstName:   shuttle.DriverFirstName,
			DriverLastName:    shuttle.DriverLastName,
			DriverGender:      shuttle.DriverGender,
			VehicleUUID:       shuttle.VehicleUUID,
			VehicleName:       shuttle.VehicleName,
			VehicleType:       shuttle.VehicleType,
			VehicleColor:      shuttle.VehicleColor,
			VehicleNumber:     shuttle.VehicleNumber,
		}
		responses = append(responses, response)
	}

	log.Println("Successfully processed shuttle responses:", responses)
	return responses, nil
}

func (s *ShuttleService) AddShuttle(req dto.ShuttleRequest, driverUUID, createdBy string) error {
	// Log: Parsing studentUUID
	studentUUID, err := uuid.Parse(req.StudentUUID)
	if err != nil {
		log.Printf("AddShuttle: Failed to parse StudentUUID - %s", req.StudentUUID)
		return err
	}
	log.Printf("AddShuttle: Parsed studentUUID - %s", studentUUID.String())

	// Log: Parsing driverUUID
	driverUUIDParsed, err := uuid.Parse(driverUUID)
	if err != nil {
		log.Printf("AddShuttle: Failed to parse driverUUID - %s", driverUUID)
		return err
	}
	log.Printf("AddShuttle: Parsed driverUUID - %s", driverUUIDParsed.String())

	// Log: Set default status if empty
	if req.Status == "" {
		req.Status = "waiting_to_be_taken_to_school"
		log.Println("AddShuttle: Set default status to 'waiting_to_be_taken_to_school'")
	}

	// Log: Create shuttle entity
	shuttle := entity.Shuttle{
		ShuttleID:   time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		ShuttleUUID: uuid.New(),
		StudentUUID: studentUUID,
		DriverUUID:  driverUUIDParsed,
		Status:      req.Status,
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
	}
	log.Printf("AddShuttle: Created shuttle entity with ShuttleID - %d", shuttle.ShuttleID)

	// Log: Attempt to save shuttle to repository
	err = s.shuttleRepository.SaveShuttle(shuttle)
	if err != nil {
		log.Println("AddShuttle: Failed to save shuttle")
		return err
	}
	log.Println("AddShuttle: Shuttle saved successfully")

	return nil
}

func (s *ShuttleService) EditShuttleStatus(shuttleUUID, status string) error {
	shuttleUUIDParsed, err := uuid.Parse(shuttleUUID)
	if err != nil {
		return err
	}

	if err := s.shuttleRepository.UpdateShuttleStatus(shuttleUUIDParsed, status); err != nil {
		return err
	}

	return nil
}
