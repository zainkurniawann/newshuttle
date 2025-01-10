package services

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"

	"github.com/google/uuid"
)

type SchoolServiceInterface interface {
	GetAllSchools(page, limit int, sortField, sortDirection string) ([]dto.SchoolResponseDTO, int, error)
	GetSpecSchool(uuid string) (dto.SchoolResponseDTO, error)
	AddSchool(req dto.SchoolRequestDTO, username string) error
	UpdateSchool(id string, req dto.SchoolRequestDTO, username string) error
	DeleteSchool(id, username, adminUUID string) error
}

type SchoolService struct {
	schoolRepository repositories.SchoolRepositoryInterface
	userRepository   repositories.UserRepositoryInterface
}

func NewSchoolService(schoolRepository repositories.SchoolRepositoryInterface, userRepository repositories.UserRepositoryInterface) SchoolService {
	return SchoolService{
		schoolRepository: schoolRepository,
		userRepository:   userRepository,
	}
}

func (service *SchoolService) GetAllSchools(page, limit int, sortField, sortDirection string) ([]dto.SchoolResponseDTO, int, error) {
	offset := (page - 1) * limit

	// Fetch data schools dan admin
	schools, adminMap, err := service.schoolRepository.FetchAllSchools(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	// Hitung total schools
	total, err := service.schoolRepository.CountSchools()
	if err != nil {
		return nil, 0, err
	}

	// Convert schools dan admin menjadi DTO
	var schoolsDTO []dto.SchoolResponseDTO
	for _, school := range schools {
		admins := adminMap[school.UUID.String()] // Ambil slice admin berdasarkan UUID sekolah

		// Buat string adminFullName (gabungkan nama admin)
		var adminFullName string
		if len(admins) == 0 {
			adminFullName = "N/A"
		} else {
			var names []string
			for _, admin := range admins {
				names = append(names, admin.FirstName+" "+admin.LastName)
			}
			adminFullName = strings.Join(names, ", ") // Gabungkan nama admin dengan koma
		}

		// Masukkan data ke DTO
		schoolsDTO = append(schoolsDTO, dto.SchoolResponseDTO{
			UUID:      school.UUID.String(),
			Name:      school.Name,
			AdminName: adminFullName,
			Address:   school.Address,
			Contact:   school.Contact,
			Email:     school.Email,
		})
	}

	return schoolsDTO, total, nil
}

func (service *SchoolService) GetSpecSchool(id string) (dto.SchoolResponseDTO, error) {
	// Mengambil data sekolah dan admin terkait
	school, admins, err := service.schoolRepository.FetchSpecSchool(id)
	if err != nil {
		return dto.SchoolResponseDTO{}, err
	}

	// Menyiapkan list nama dan UUID admin
	var adminNames, adminUUIDs []string
	for _, admin := range admins {
		userFullName := "N/A"
		if admin.UserUUID != uuid.Nil {
			userFullName = admin.FirstName + " " + admin.LastName
		}

		adminUUID := "N/A"
		if admin.UserUUID != uuid.Nil {
			adminUUID = admin.UserUUID.String()
		}

		// Menggabungkan nama dan UUID admin
		adminNames = append(adminNames, userFullName)
		adminUUIDs = append(adminUUIDs, adminUUID)
	}

	// Menggabungkan list nama dan UUID admin menjadi string
	adminUUIDsStr := strings.Join(adminUUIDs, ", ")
	adminNamesStr := strings.Join(adminNames, ", ")

	// Memeriksa apakah school.Point valid
	var pointStr string
	if school.Point.Valid {
		pointStr = school.Point.String
	}

	// Menyiapkan SchoolResponseDTO
	schoolDTO := dto.SchoolResponseDTO{
		UUID:        school.UUID.String(),
		Name:        school.Name,
		AdminUUID:   adminUUIDsStr,
		AdminName:   adminNamesStr,
		Address:     school.Address,
		Contact:     school.Contact,
		Email:       school.Email,
		Description: school.Description,
		Point:       pointStr, // Menambahkan point jika valid
		CreatedAt:   safeTimeFormat(school.CreatedAt),
		CreatedBy:   safeStringFormat(school.CreatedBy),
		UpdatedAt:   safeTimeFormat(school.UpdatedAt),
		UpdatedBy:   safeStringFormat(school.UpdatedBy),
	}

	return schoolDTO, nil
}


func (service *SchoolService) AddSchool(req dto.SchoolRequestDTO, username string) error {
	// Convert map to JSON string (handling empty Point)
	var pointJSON string
	if len(req.Point) > 0 {
		pointData, err := json.Marshal(req.Point)
		if err != nil {
			return err
		}
		pointJSON = string(pointData)
	} else {
		// If Point is empty, set it to an empty JSON object or handle as needed
		pointJSON = "{}"
	}

	school := entity.School{
		ID:          time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		UUID:        uuid.New(),
		Name:        req.Name,
		Address:     req.Address,
		Contact:     req.Contact,
		Email:       req.Email,
		Description: req.Description,
		// Save the Point as a non-null value (if empty, use "{}" or other default value)
		Point:       sql.NullString{String: pointJSON, Valid: true}, // Make sure Valid is true
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		CreatedBy:   toNullString(username),
	}

	if err := service.schoolRepository.SaveSchool(school); err != nil {
		return err
	}

	return nil
}




func (service *SchoolService) UpdateSchool(id string, req dto.SchoolRequestDTO, username string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	pointJSON, err := json.Marshal(req.Point)
	if err != nil {
		return err
	}

	school := entity.School{
		UUID:        parsedUUID,
		Name:        req.Name,
		Address:     req.Address,
		Contact:     req.Contact,
		Email:       req.Email,
		Description: req.Description,
		Point:       sql.NullString{String: string(pointJSON), Valid: true},
		UpdatedAt:   toNullTime(time.Now()),
		UpdatedBy:   toNullString(username),
	}

	if err := service.schoolRepository.UpdateSchool(school); err != nil {
		return err
	}

	return nil
}

func (service *SchoolService) DeleteSchool(id, username, adminUUID string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Handle multiple admin UUIDs deletion
	if adminUUID != "N/A" && adminUUID != "" {
		uuidList := strings.Split(adminUUID, ", ")

		tx, err := service.userRepository.BeginTransaction()
		if err != nil {
			return err
		}
		var transactionErr error
		defer func() {
			if transactionErr != nil {
				tx.Rollback()
			} else {
				transactionErr = tx.Commit()
			}
		}()

		for _, uuids := range uuidList {
			parsedAdminUUID, err := uuid.Parse(uuids)
			if err != nil {
				continue
			}

			if err := service.userRepository.DeleteSchoolAdmin(tx, parsedAdminUUID, username); err != nil {
				return err
			}
		}

		school := entity.School{
			UUID:      parsedUUID,
			DeletedAt: toNullTime(time.Now()),
			DeletedBy: toNullString(username),
		}

		if err := service.schoolRepository.DeleteSchool(school); err != nil {
			return err
		}

		return nil
	} else {
		tx, err := service.userRepository.BeginTransaction()
		if err != nil {
			return err
		}
		var transactionErr error
		defer func() {
			if transactionErr != nil {
				tx.Rollback()
			} else {
				transactionErr = tx.Commit()
			}
		}()

		// Delete the school
		school := entity.School{
			UUID:      parsedUUID,
			DeletedAt: toNullTime(time.Now()),
			DeletedBy: toNullString(username),
		}
		if err := service.schoolRepository.DeleteSchool(school); err != nil {
			return err
		}

		return nil
	}
}

func safeStringFormat(s sql.NullString) string {
	if !s.Valid || s.String == "" {
		return "N/A"
	}
	return s.String
}

func safeTimeFormat(t sql.NullTime) string {
	if !t.Valid {
		return "N/A"
	}
	return t.Time.Format(time.RFC3339)
}

func toNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{String: value, Valid: false}
	}
	return sql.NullString{String: value, Valid: true}
}

func toNullTime(value time.Time) sql.NullTime {
	return sql.NullTime{Time: value, Valid: !value.IsZero()}
}
