package services

import (
	"database/sql"
	"fmt"
	"log"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type RouteServiceInterface interface {
	GetAllRoutesByAS(page, limit int, sortField, sortDirection, schoolUUID string) ([]dto.RoutesResponseDTO, int, error)
	GetSpecRouteByAS(routeNameUUID, driverUUID string) (dto.RoutesResponseDTO, error)
	GetAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)
	AddRoute(route dto.RoutesRequestDTO, schoolUUID, username string) error
	UpdateRoute(request dto.UpdateRouteRequest, routeNameUUID, schoolUUID, username string) error
	GetDriverUUIDByRouteName(routeNameUUID string) (string, error)
	DeleteRoute(routenameUUID, schoolUUID, username string) error
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
		studentOrder, err := strconv.Atoi(route.StudentOrder)
    if err != nil {
        fmt.Printf("Error parsing StudentOrder for %s: %v\n", route.StudentUUID, err)
        continue // Skip data yang bermasalah
    }
		student := dto.StudentDTO{
			RouteAssignmentUUID: route.RouteAssignmentUUID.String(),
			StudentUUID:      route.StudentUUID.String(),
			StudentFirstName: defaultString(route.StudentFirstName),
			StudentLastName:  defaultString(route.StudentLastName),
			StudentStatus: route.StudentStatus,
			StudentOrder:     studentOrder,
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
	log.Println("[AddRoute] Mulai menambahkan rute baru.")

	// Validasi duplicate student
	if err := ValidateDuplicateStudents(route.RouteAssignment); err != nil {
		log.Printf("[AddRoute] Validasi duplicate students gagal: %v\n", err)
		return err
	}

	// Buat entitas route
	routeEntity := entity.Routes{
		RouteID:          time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		RouteNameUUID:    uuid.New(),
		SchoolUUID:       uuid.MustParse(schoolUUID),
		RouteName:        route.RouteName,
		RouteDescription: route.RouteDescription,
		CreatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
		CreatedBy:        sql.NullString{String: username, Valid: true},
	}
	log.Printf("[AddRoute] Route entity dibuat: %+v\n", routeEntity)

	// Mulai transaksi database
	tx, err := service.routeRepository.BeginTransaction()
	if err != nil {
		log.Printf("[AddRoute] Gagal memulai transaksi: %v\n", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	log.Println("[AddRoute] Transaksi database dimulai.")

	// Pastikan transaksi dibatalkan jika ada panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Println("[AddRoute] Transaksi dibatalkan karena panic.")
			panic(r)
		}
	}()

	// Insert route ke database
	routeNameUUID, err := service.routeRepository.AddRoutes(tx, routeEntity)
	if err != nil {
		log.Printf("[AddRoute] Gagal menambahkan rute: %v\n", err)
		tx.Rollback()
		return fmt.Errorf("failed to add route: %w", err)
	}
	log.Printf("[AddRoute] Rute berhasil ditambahkan dengan routeNameUUID: %s\n", routeNameUUID)

	parsedRouteUUID := uuid.MustParse(routeNameUUID)

	// Proses penugasan driver dan siswa
	for _, assignment := range route.RouteAssignment {
		log.Printf("[AddRoute] Memproses assignment untuk driver: %s\n", assignment.DriverUUID.String())
		isDriverAssigned, err := service.routeRepository.IsDriverAssigned(tx, assignment.DriverUUID.String())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error checking driver assignment: %w", err)
		}
		if isDriverAssigned {
			tx.Rollback()
			return fmt.Errorf("driver already assigned to another route")
		}

		// Validasi kapasitas kursi kendaraan
		vehicleSeats, err := service.routeRepository.GetVehicleSeatsByDriver(tx, assignment.DriverUUID.String())
		if err != nil {
			log.Printf("[AddRoute] Gagal mendapatkan kapasitas kursi kendaraan: %v\n", err)
			tx.Rollback()
			return fmt.Errorf("error fetching vehicle seats: %w", err)
		}
		log.Printf("[AddRoute] Kapasitas kursi kendaraan untuk driver %s: %d\n", assignment.DriverUUID.String(), vehicleSeats)

		assignedStudentsCount, err := service.routeRepository.CountAssignedStudentsByDriver(tx, assignment.DriverUUID.String())
		if err != nil {
			log.Printf("[AddRoute] Gagal menghitung jumlah siswa yang diassign ke driver: %v\n", err)
			tx.Rollback()
			return fmt.Errorf("error fetching assigned students count: %w", err)
		}
		log.Printf("[AddRoute] Jumlah siswa yang sudah diassign untuk driver %s: %d\n", assignment.DriverUUID.String(), assignedStudentsCount)

		// Validasi kapasitas kursi
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
			if assignedStudentsCount+1 > vehicleSeats {
				log.Printf("[AddRoute] Kapasitas kursi melebihi batas: %d/%d\n", assignedStudentsCount, vehicleSeats)
				tx.Rollback()
				return fmt.Errorf("Maximum seats exceeded, Capacity is %d", vehicleSeats)
			}						

			if student.StudentOrder == "" || student.StudentOrder == "0" {
				log.Println("[AddRoute] Student order kosong atau nol.")
				tx.Rollback()
				return fmt.Errorf("student order cannot be empty or zero")
			}

			routeAssignmentEntity := entity.RouteAssignment{
				RouteID:       			time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
				RouteAssignmentUUID:    uuid.New(),
				DriverUUID:    			assignment.DriverUUID,
				StudentUUID:   			student.StudentUUID,
				StudentOrder:  			student.StudentOrder,
				SchoolUUID:    			uuid.MustParse(schoolUUID),
				RouteNameUUID: 			parsedRouteUUID.String(),
				CreatedAt:     			sql.NullTime{Time: time.Now(), Valid: true},
				CreatedBy:     			sql.NullString{String: username, Valid: true},
			}
			log.Printf("[AddRoute] Menambahkan route assignment: %+v\n", routeAssignmentEntity)

			if err := service.routeRepository.AddRouteAssignment(tx, routeAssignmentEntity); err != nil {
				log.Printf("[AddRoute] Gagal menambahkan route assignment: %v\n", err)
				tx.Rollback()
				return fmt.Errorf("failed to add route assignment: %w", err)
			}

			// Update jumlah siswa yang diassign
			assignedStudentsCount++
			log.Printf("[AddRoute] Jumlah siswa yang diassign diperbarui: %d\n", assignedStudentsCount)
		}
	}

	// Commit transaksi
	if err := tx.Commit(); err != nil {
		log.Printf("[AddRoute] Gagal commit transaksi: %v\n", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	log.Println("[AddRoute] Transaksi berhasil di-commit.")

	log.Println("[AddRoute] Rute berhasil ditambahkan.")
	return nil
}

func (s *routeService) UpdateRoute(requestDTO dto.UpdateRouteRequest, routeNameUUID, schoolUUID, username string) error {
    // Validasi input
     // Validasi input
	 if routeNameUUID == "" || requestDTO.DriverUUID == "" {
        return fmt.Errorf("RouteNameUUID dan DriverUUID tidak boleh kosong")
    }

    // Pastikan UUID valid
    parsedRouteUUID := uuid.MustParse(routeNameUUID)
    parsedSchoolUUID := uuid.MustParse(schoolUUID) // Pastikan SchoolUUID juga diparse dengan benar

    // Update data route
    routeEntity := entity.Routes{
        RouteNameUUID:    parsedRouteUUID, // Menambahkan RouteNameUUID
        SchoolUUID:       parsedSchoolUUID, // Menambahkan SchoolUUID
        RouteName:        requestDTO.RouteName,
        RouteDescription: requestDTO.RouteDescription,
        UpdatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
        UpdatedBy:        sql.NullString{String: username, Valid: true},
    }
    
    if err := s.routeRepository.UpdateRouteDetails(&routeEntity); err != nil {
        return fmt.Errorf("Gagal memperbarui data route: %w", err)
    }

    // Menambahkan siswa baru
    for _, student := range requestDTO.Added {
        // Validasi StudentOrder
        if student.StudentOrder <= 0 {
            log.Printf("Invalid StudentOrder for student %s. Using default value 1.", student.StudentUUID)
            student.StudentOrder = 1 // Gunakan nilai default
        }
		studentUUID, err := uuid.Parse(student.StudentUUID)
        if err != nil {
            return fmt.Errorf("StudentUUID tidak valid: %w", err)
        }
		driverUUID, err := uuid.Parse(requestDTO.DriverUUID)
        if err != nil {
            return fmt.Errorf("StudentUUID tidak valid: %w", err)
        }
		studentOrder := strconv.Itoa(student.StudentOrder)
			if err != nil {
				return fmt.Errorf("schoolUUID tidak valid: %w", err)
			}

        assignmentEntity := entity.RouteAssignment{
			RouteID:       			time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
			RouteAssignmentUUID:    uuid.New(),
			DriverUUID:    			driverUUID,
			StudentUUID:   			studentUUID,
			StudentOrder:  			studentOrder,
			SchoolUUID:    			uuid.MustParse(schoolUUID),
			RouteNameUUID: 			parsedRouteUUID.String(),
			CreatedAt:     			sql.NullTime{Time: time.Now(), Valid: true},
			CreatedBy:     			sql.NullString{String: username, Valid: true},
        }
        if err := s.routeRepository.AddStudentToRoute(&assignmentEntity); err != nil {
            return fmt.Errorf("Gagal menambahkan siswa ke route: %w", err)
        }
    }

    // Menghapus siswa
    for _, student := range requestDTO.DeletedStudents {
        if err := s.routeRepository.DeleteStudentFromRoute(routeNameUUID, student.StudentUUID, schoolUUID); err != nil {
            return fmt.Errorf("Gagal menghapus siswa dari route: %w", err)
        }
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