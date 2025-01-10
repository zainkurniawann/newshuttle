package services

import (
	"database/sql"
	"fmt"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"
	"github.com/google/uuid"
)

type RouteServiceInterface interface {
	GetAllRoutesByAS(page, limit int, sortField, sortDirection, schoolUUID string) ([]dto.RoutesResponseDTO, int, error)
	GetSpecRouteByAS(routeNameUUID, driverUUID string) (dto.RoutesResponseDTO, error)
	GetAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)

	AddRoute(route dto.RoutesRequestDTO, schoolUUID, username string) error
	UpdateRoute(route dto.RoutesRequestDTO, routenameUUID, schoolUUID, username string) error 
	DeleteRoute(routenameUUID, schoolUUID, username string) error

	GetDriverUUIDByRouteName(routeNameUUID string) (string, error)
}

type routeService struct {
	routeRepository repositories.RouteRepositoryInterface
}

func NewRouteService(routeRepository repositories.RouteRepositoryInterface) RouteServiceInterface {
	return &routeService{
		routeRepository: routeRepository,
	}
}

func (service *routeService) GetAllRoutesByAS(page, limit int, sortField, sortDirection, schoolUUID string) ([]dto.RoutesResponseDTO, int, error) {
	offset := (page - 1) * limit

	// Panggil repository untuk mendapatkan data dan total items
	routes, err := service.routeRepository.FetchAllRoutesByAS(offset, limit, sortField, sortDirection, schoolUUID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get routes: %w", err)
	}

	totalItems, err := service.routeRepository.CountRoutesBySchool(schoolUUID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count routes: %w", err)
	}

	return routes, totalItems, nil
}


func (s *routeService) GetSpecRouteByAS(routeNameUUID, driverUUID string) (dto.RoutesResponseDTO, error) {
	if driverUUID == "" {
		driverUUID = ""
	}

	routes, err := s.routeRepository.FetchSpecRouteByAS(routeNameUUID, driverUUID)
	if err != nil {
		return dto.RoutesResponseDTO{}, err
	}

	var routeResponse dto.RoutesResponseDTO
	if len(routes) == 0 {
		routeResponse.RouteName = "Route not assigned"
		routeResponse.RouteDescription = "No description available"
		routeResponse.RouteAssignment = nil
		return routeResponse, nil
	}

	routeResponse.RouteNameUUID = routes[0].RouteNameUUID
	routeResponse.RouteName = routes[0].RouteName
	routeResponse.RouteDescription = routes[0].RouteDescription

	if routes[0].DriverUUID == uuid.Nil {
		routeResponse.RouteAssignment = nil
		return routeResponse, nil
	}
	driverInfo := dto.RouteAssignmentResponseDTO{
		DriverUUID:      routes[0].DriverUUID.String(),
		DriverFirstName: defaultString(routes[0].DriverFirstName),
		DriverLastName:  defaultString(routes[0].DriverLastName),
	}

	for _, route := range routes {
		student := dto.StudentDTO{
			StudentUUID:      route.StudentUUID.String(),
			StudentFirstName: defaultString(route.StudentFirstName),
			StudentLastName:  defaultString(route.StudentLastName),
			StudentStatus: route.StudentStatus,
			StudentOrder:     route.StudentOrder,
		}
		driverInfo.Students = append(driverInfo.Students, student)
	}

	routeResponse.RouteAssignment = append(routeResponse.RouteAssignment, driverInfo)
	return routeResponse, nil
}

func defaultString(str string) string {
	if str == "" {
		return "Unknown"
	}
	return str
}

func (service *routeService) GetAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error) {
	routes, err := service.routeRepository.FetchAllRoutesByDriver(driverUUID)
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (service *routeService) AddRoute(route dto.RoutesRequestDTO, schoolUUID, username string) error {

	if err := ValidateDuplicateStudents(route.RouteAssignment); err != nil {
		return err
	}

	routeEntity := entity.Routes{
		RouteID:          time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		RouteNameUUID:    uuid.New(),
		SchoolUUID:       uuid.MustParse(schoolUUID),
		RouteName:        route.RouteName,
		RouteDescription: route.RouteDescription,
		CreatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
		CreatedBy:        sql.NullString{String: username, Valid: true},
	}

	tx, err := service.routeRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Insert route ke database
	routeNameUUID, err := service.routeRepository.AddRoutes(tx, routeEntity)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to add route: %w", err)
	}

	parsedRouteUUID := uuid.MustParse(routeNameUUID)

	// Insert route assignment
	for _, assignment := range route.RouteAssignment {
		isDriverAssigned, err := service.routeRepository.IsDriverAssigned(tx, assignment.DriverUUID.String())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error checking driver assignment: %w", err)
		}
		if isDriverAssigned {
			tx.Rollback()
			return fmt.Errorf("driver already assigned to another route")
		}

		for _, student := range assignment.Students {
			isStudentAssigned, err := service.routeRepository.IsStudentAssigned(tx, student.StudentUUID.String())
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error checking student assignment: %w", err)
			}
			if isStudentAssigned {
				tx.Rollback()
				return fmt.Errorf("student already assigned to another route")
			}

			if student.StudentOrder == "" || student.StudentOrder == "0" {
				tx.Rollback()
				return fmt.Errorf("student order cannot be empty or zero")
			}

			routeAssignmentEntity := entity.RouteAssignment{
				RouteID:       time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
				RouteUUID:     parsedRouteUUID,
				DriverUUID:    assignment.DriverUUID,
				StudentUUID:   student.StudentUUID,
				StudentOrder:  student.StudentOrder,
				SchoolUUID:    uuid.MustParse(schoolUUID),
				RouteNameUUID: routeEntity.RouteNameUUID.String(),
				CreatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
				CreatedBy:     sql.NullString{String: username, Valid: true},
			}

			if err := service.routeRepository.AddRouteAssignment(tx, routeAssignmentEntity); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to add route assignment: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (service *routeService) UpdateRoute(route dto.RoutesRequestDTO, routenameUUID, schoolUUID, username string) error {
	routeEntity := entity.Routes{
		RouteNameUUID:    uuid.MustParse(routenameUUID),
		SchoolUUID:       uuid.MustParse(schoolUUID),
		RouteName:        route.RouteName,
		RouteDescription: route.RouteDescription,
		UpdatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:        sql.NullString{String: username, Valid: true},
	}

	tx, err := service.routeRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err = service.routeRepository.UpdateRoute(tx, routeEntity)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update route: %w", err)
	}

	for _, assignment := range route.RouteAssignment {
		for _, student := range assignment.Students {
			routeAssignmentEntity := entity.RouteAssignment{
				RouteUUID:     uuid.MustParse(routenameUUID),
				DriverUUID:    assignment.DriverUUID,
				StudentUUID:   student.StudentUUID,
				StudentOrder:  student.StudentOrder,
				SchoolUUID:    uuid.MustParse(schoolUUID),
				CreatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
				CreatedBy:     sql.NullString{String: username, Valid: true},
			}

			err := service.routeRepository.UpdateOrAddRouteAssignment(tx, routeAssignmentEntity)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update route assignment: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (service *routeService) DeleteRoute(routenameUUID, schoolUUID, username string) error {
	tx, err := service.routeRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	routeExists, err := service.routeRepository.RouteExists(tx, routenameUUID, schoolUUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error checking if route exists: %w", err)
	}
	if !routeExists {
		tx.Rollback()
		return fmt.Errorf("route not found")
	}

	err = service.routeRepository.DeleteRouteAssignments(tx, routenameUUID, schoolUUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting route assignments: %w", err)
	}

	err = service.routeRepository.DeleteRoute(tx, routenameUUID, schoolUUID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting route: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *routeService) GetDriverUUIDByRouteName(routeNameUUID string) (string, error) {
	driverUUID, err := s.routeRepository.GetDriverUUIDByRouteName(routeNameUUID)
	if err != nil {
		return "", fmt.Errorf("error retrieving driver UUID: %w", err)
	}
	return driverUUID, nil
}

func ValidateDuplicateStudents(routeAssignments []dto.RouteAssignmentRequestDTO) error {
	studentSet := make(map[string]bool)

	for _, assignment := range routeAssignments {
		for _, student := range assignment.Students {
			if studentSet[student.StudentUUID.String()] {
				// Kembalikan error dengan pesan relevan
				return fmt.Errorf("same student not permitted")
			}
			studentSet[student.StudentUUID.String()] = true
		}
	}

	return nil
}
