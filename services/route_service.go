package services

import (
	"database/sql"
	"fmt"
	"log"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RouteServiceInterface interface {
	GetAllRoutesByAS(page, limit int, sortField, sortDirection, schoolUUID string) ([]dto.RoutesResponseDTO, int, error)
	GetAllRouteAssignments(page, limit int, sortField, sortDirection string) ([]dto.RoutesResponseDTO, int, error)
	GetSpecRouteByAS(routeNameUUID, driverUUID string) (dto.RoutesResponseDTO, error)
	GetAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)
	AddRoute(route dto.RoutesRequestDTO, schoolUUID, username string) error
	UpdateRoute(request dto.UpdateRouteRequest, routeNameUUID, schoolUUID, username string) error
	UpdateStudentOrderByDriver(studentUUID string, newOrder int) error
	GetMaxStudentOrder(routeNameUUID, schoolUUID string) (int, error)
	GetDriverUUIDByRouteName(routeNameUUID string) (string, error)
	DeleteRoute(routenameUUID, schoolUUID, username string) error

	GetTotalDistance(driverStart [2]float64, students [][2]float64, school [2]float64) float64
}

type routeService struct {
	routeRepository repositories.RouteRepositoryInterface
}

func NewRouteService(routeRepository repositories.RouteRepositoryInterface) RouteServiceInterface {
	return &routeService{
		routeRepository: routeRepository,
	}
}

func (s *routeService) GetTotalDistance(driverStart [2]float64, students [][2]float64, school [2]float64) float64 {
	return s.routeRepository.CalculateTotalDistance(driverStart, students, school)
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

func (s *routeService) GetAllRouteAssignments(page, limit int, sortField, sortDirection string) ([]dto.RoutesResponseDTO, int, error) {
    log.Println("Starting GetAllRouteAssignments service")

    // Ambil data dari repository
    log.Printf("Fetching data from repository with page: %d, limit: %d\n", page, limit)
    routes, totalItems, err := s.routeRepository.FetchAllRouteAssignments(page, limit)
    if err != nil {
        log.Printf("Failed to fetch data from repository: %v\n", err)
        return nil, 0, err
    }
    log.Printf("Fetched %d items. Total items available: %d\n", len(routes), totalItems)

    // Terapkan logika pengurutan
    log.Printf("Applying sorting on field: %s with direction: %s\n", sortField, sortDirection)
    sort.Slice(routes, func(i, j int) bool {
        switch strings.ToLower(sortField) {
        case "route_name":
            if strings.ToLower(sortDirection) == "desc" {
                return routes[i].RouteName > routes[j].RouteName
            }
            return routes[i].RouteName < routes[j].RouteName
        case "driver_first_name":
            if strings.ToLower(sortDirection) == "desc" {
                return routes[i].RouteAssignment[0].DriverFirstName > routes[j].RouteAssignment[0].DriverFirstName
            }
            return routes[i].RouteAssignment[0].DriverFirstName < routes[j].RouteAssignment[0].DriverFirstName
        default:
            // Default sorting berdasarkan RouteNameUUID
            if strings.ToLower(sortDirection) == "desc" {
                return routes[i].RouteNameUUID > routes[j].RouteNameUUID
            }
            return routes[i].RouteNameUUID < routes[j].RouteNameUUID
        }
    })
    log.Println("Sorting applied successfully")

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
	log.Println("Getting all routes for driver:", driverUUID)
	routes, err := service.routeRepository.FetchAllRoutesByDriver(driverUUID)
	if err != nil {
		log.Println("Error in routeRepository:", err)
		return nil, err
	}
	log.Println("Routes retrieved from repository:", len(routes))
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
    if routeNameUUID == "" || requestDTO.DriverUUID == "" {
        return fmt.Errorf("RouteNameUUID dan DriverUUID tidak boleh kosong")
    }

    // Pastikan UUID valid
    parsedRouteUUID := uuid.MustParse(routeNameUUID)
    parsedSchoolUUID := uuid.MustParse(schoolUUID)

    // Update data route
    routeEntity := entity.Routes{
        RouteNameUUID:    parsedRouteUUID,
        SchoolUUID:       parsedSchoolUUID,
        RouteName:        requestDTO.RouteName,
        RouteDescription: requestDTO.RouteDescription,
        UpdatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
        UpdatedBy:        sql.NullString{String: username, Valid: true},
    }

    if err := s.routeRepository.UpdateRouteDetails(&routeEntity); err != nil {
        return fmt.Errorf("Gagal memperbarui data route: %w", err)
    }

    // Mendapatkan StudentOrder terbesar
    maxStudentOrder, err := s.GetMaxStudentOrder(routeNameUUID, schoolUUID)
    if err != nil {
        return fmt.Errorf("Gagal mendapatkan StudentOrder terbesar: %w", err)
    }

    // Menambahkan siswa baru
    for _, student := range requestDTO.Added {
        if student.StudentOrder <= 0 {
            student.StudentOrder = maxStudentOrder + 1
            maxStudentOrder++
        }

        studentUUID, err := uuid.Parse(student.StudentUUID)
        if err != nil {
            return fmt.Errorf("StudentUUID tidak valid: %w", err)
        }

        driverUUID, err := uuid.Parse(requestDTO.DriverUUID)
        if err != nil {
            return fmt.Errorf("DriverUUID tidak valid: %w", err)
        }

        assignmentEntity := entity.RouteAssignment{
            RouteID:              time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
            RouteAssignmentUUID:  uuid.New(),
            DriverUUID:           driverUUID,
            StudentUUID:          studentUUID,
            StudentOrder:         strconv.Itoa(student.StudentOrder),
            SchoolUUID:           parsedSchoolUUID,
            RouteNameUUID:        parsedRouteUUID.String(),
            CreatedAt:            sql.NullTime{Time: time.Now(), Valid: true},
            CreatedBy:            sql.NullString{String: username, Valid: true},
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

    // Mengupdate urutan siswa berdasarkan drag-and-drop
    for _, student := range requestDTO.Students {
		log.Printf("Processing update for student %s with StudentOrder %d\n", student.StudentUUID, student.StudentOrder)
	
		assignmentEntity := entity.RouteAssignment{
			RouteNameUUID: parsedRouteUUID.String(),  // Gunakan RouteNameUUID
			StudentOrder:  strconv.Itoa(student.StudentOrder),      // Pastikan ini adalah nilai yang benar
			UpdatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedBy:     sql.NullString{String: username, Valid: true},
		}
	
		log.Printf("Assignment entity created: %+v\n", assignmentEntity)
	
		// Update student order berdasarkan StudentUUID
		if err := s.routeRepository.UpdateStudentOrder(routeNameUUID, &assignmentEntity, student.StudentUUID); err != nil {
			log.Printf("Error updating student order for student with StudentUUID %s: %v\n", student.StudentUUID, err)
			return fmt.Errorf("Gagal memperbarui urutan siswa: %w", err)
		}
	
		log.Printf("Successfully updated student order for student with StudentUUID %s\n", student.StudentUUID)
	}
			

    return nil
}

func (service *routeService) UpdateStudentOrderByDriver(studentUUID string, newOrder int) error {
	err := service.routeRepository.UpdateStudentOrderByDriver(studentUUID, newOrder)
	if err != nil {
		log.Println("Error in routeRepository:", err)
		return err
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

func (s *routeService) GetMaxStudentOrder(routeNameUUID, schoolUUID string) (int, error) {
    // Panggil repository untuk mendapatkan nilai StudentOrder terbesar
    maxOrder, err := s.routeRepository.GetMaxStudentOrder(routeNameUUID, schoolUUID)
    if err != nil {
        return 0, fmt.Errorf("Gagal mendapatkan StudentOrder terbesar: %w", err)
    }
    return maxOrder, nil
}