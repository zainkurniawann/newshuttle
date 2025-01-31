package services

// import (
// 	// "database/sql"
// 	// "encoding/json"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"

// 	// "log"

// 	"shuttle/errors"
// 	// "shuttle/logger"
// 	"shuttle/models/dto"
// 	"shuttle/models/entity"
// 	"shuttle/repositories"
// 	"time"

// 	"github.com/google/uuid"
// 	"github.com/jmoiron/sqlx"
// )

// type RegisterServiceInterface interface {
// 	AddUserRegister(req dto.UserRequestsDTO, user_name string) (uuid.UUID, error)
// 	ApproveUserRegister(userUUID uuid.UUID, username string, approverRole string) error
// }

// type RegisterService struct {
// 	registerRepository repositories.RegisterRepositoryInterface
// }

// func NewRegisterService(registerRepository repositories.RegisterRepositoryInterface) RegisterService {
// 	return RegisterService{
// 		registerRepository: registerRepository,
// 	}
// }

// func (s *RegisterService) AddUserRegister(req dto.UserRequestsDTO, user_name string) (uuid.UUID, error) {
// 	if req.RoleCode == "SA" {
// 		return uuid.Nil, fmt.Errorf("superadmin tidak dapat melakukan self-registration")
// 	}

// 	exists, err := s.registerRepository.CheckEmailExistForRegister("", req.Email)
// 	if err != nil {
// 		return uuid.Nil, err
// 	}
// 	if exists {
// 		return uuid.Nil, fmt.Errorf("email %s sudah terdaftar", req.Email)
// 	}

// 	exists, err = s.registerRepository.CheckUsernameExistForRegister("", req.Username)
// 	if err != nil {
// 		return uuid.Nil, err
// 	}
// 	if exists {
// 		return uuid.Nil, fmt.Errorf("username %s sudah terdaftar", req.Username)
// 	}

// 	if req.Password != "" {
// 		hashedPassword, err := hashPassword(req.Password)
// 		if err != nil {
// 			return uuid.Nil, err
// 		}
// 		req.Password = hashedPassword
// 	}

// 	tx, err := s.registerRepository.BeginTransaction()
// 	if err != nil {
// 		return uuid.Nil, fmt.Errorf("gagal memulai transaksi: %w", err)
// 	}

// 	var transactionErr error
// 	defer func() {
// 		if transactionErr != nil {
// 			tx.Rollback()
// 		} else {
// 			transactionErr = tx.Commit()
// 		}
// 	}()

// 	userEntity := entity.User{
// 		ID:             time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
// 		UUID:           uuid.New(),
// 		Username:       req.Username,
// 		Email:          req.Email,
// 		Password:       req.Password,
// 		Role:           entity.Role(req.Role),
// 		RoleCode:       req.RoleCode,
// 		RegisterStatus: "PENDING", // Selalu pending untuk approval
// 	}

// 	userUUID, err := s.registerRepository.SaveUserRegister(tx, userEntity)
// 	if err != nil {
// 		transactionErr = fmt.Errorf("gagal menyimpan user: %w", err)
// 		return uuid.Nil, transactionErr
// 	}

// 	if err := s.saveRoleRegisterDetails(tx, userEntity.UUID, req); err != nil {
// 		transactionErr = fmt.Errorf("gagal menyimpan detail user: %w", err)
// 		return uuid.Nil, transactionErr
// 	}

// 	return userUUID, nil
// }

// func (s *RegisterService) saveRoleRegisterDetails(tx *sqlx.Tx, userUUID uuid.UUID, req dto.UserRequestsDTO) error {
// 	switch entity.Role(req.Role) {
// 	case entity.SchoolAdmin:
// 		parsedDetails, err := parseDetails[dto.SchoolAdminDetailsRequestsDTO](req.Details)
// 		if err != nil {
// 			return errors.New("invalid school admin details format: " + err.Error(), 400)
// 		}
// 		if parsedDetails.SchoolUUID == "" && (parsedDetails.School.Name == "" || parsedDetails.School.Address == "") {
// 			return errors.New("data sekolah tidak lengkap untuk membuat sekolah baru", 400)
// 		}
// 		var schoolUUID uuid.UUID
// 		if parsedDetails.SchoolUUID == "" || parsedDetails.SchoolUUID == "nil" {
// 			schoolUUID, err = s.createSchoolForRegister(tx, parsedDetails.School)
// 			if err != nil {
// 				return fmt.Errorf("gagal membuat sekolah: %w", err)
// 			}
// 		} else {
// 			schoolUUID = uuid.MustParse(parsedDetails.SchoolUUID)
// 		}
// 		schoolDetails := entity.SchoolAdminDetails{
// 			SchoolUUID: schoolUUID,
// 			Picture:    req.Picture,
// 			FirstName:  req.FirstName,
// 			LastName:   req.LastName,
// 			Gender:     entity.Gender(req.Gender),
// 			Phone:      req.Phone,
// 			Address:    req.Address,
// 		}
// 		return s.registerRepository.SaveSchoolAdminRegisterDetails(tx, schoolDetails, userUUID, nil)

// 	case entity.Parent:
// 		details := entity.ParentDetails{
// 			FirstName: req.FirstName,
// 			LastName:  req.LastName,
// 			Gender:    entity.Gender(req.Gender),
// 			Phone:     req.Phone,
// 			Address:   req.Address,
// 		}
// 		return s.registerRepository.SaveParentRegisterDetails(tx, details, userUUID, nil)

// 	case entity.Driver:
// 		parsedDetails, err := parseDetails[dto.DriverDetailsRequestsDTO](req.Details)
// 		if err != nil {
// 			return errors.New("invalid driver details format: "+err.Error(), 400)
// 		}
		
// 		fmt.Printf("driver details: %+v\n", parsedDetails)
// 		println("vehicleUUID: ", parsedDetails.VehicleUUID)
// 		driverDetails := entity.DriverDetails{
// 			SchoolUUID:    parseSafeUUID(parsedDetails.SchoolUUID),
// 			VehicleUUID:   parseSafeUUID(parsedDetails.VehicleUUID),
// 			Picture:       req.Picture,
// 			FirstName:     req.FirstName,
// 			LastName:      req.LastName,
// 			Gender:        entity.Gender(req.Gender),
// 			Phone:         req.Phone,
// 			Address:       req.Address,
// 			LicenseNumber: parsedDetails.LicenseNumber,
// 		}
// 		return s.registerRepository.SaveDriverRegisterDetails(tx, driverDetails, userUUID, nil)

// 	default:
// 		return errors.New("invalid role", 400)
// 	}
// }

// func (s *RegisterService) ApproveUserRegister(userUUID uuid.UUID, username string, approverRole string) error {
// 	var requiredApproverRole string

// 	user, err := s.registerRepository.GetUserByUUID(userUUID)
// 	if err != nil {
// 		return fmt.Errorf("user tidak ditemukan: %w", err)
// 	}

// 	// Menentukan siapa yang berhak approve berdasarkan role yang didaftarkan
// 	switch user.RoleCode {
// 	case "P":
// 		requiredApproverRole = "AS"
// 	case "AS":
// 		requiredApproverRole = "SA"
// 	case "D":
// 		requiredApproverRole = "SA"
// 	default:
// 		return fmt.Errorf("role tidak valid untuk approval")
// 	}

// 	// Validasi role approver
// 	if approverRole != requiredApproverRole {
// 		return fmt.Errorf("anda tidak memiliki izin untuk melakukan approval")
// 	}
	
// 	err = s.registerRepository.ApproveUserRegister(username, userUUID)
// 	if err != nil {
// 		return fmt.Errorf("gagal memperbarui status user: %w", err)
// 	}

// 	return nil
// }

// func (service *RegisterService) createSchoolForRegister(tx *sqlx.Tx, req dto.SchoolRequestDTO) (uuid.UUID, error) {
// 	// Convert Point map to JSON string (handling empty Point)
// 	var point interface{}
// 	if len(req.Point) > 0 {
// 		// Marshal Point to JSON string
// 		pointData, err := json.Marshal(req.Point)
// 		if err != nil {
// 			return uuid.Nil, fmt.Errorf("failed to marshal Point: %w", err)
// 		}
// 		point = string(pointData) // Store as JSON string
// 	} else {
// 		// If Point is empty, set it to NULL
// 		point = nil
// 	}

// 	// Create School entity
// 	school := entity.School{
// 		ID:          time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6), // Custom ID
// 		UUID:        uuid.New(),
// 		Name:        req.Name,
// 		Address:     req.Address,
// 		Contact:     req.Contact,
// 		Email:       req.Email,
// 		Description: req.Description,
// 		Point:       sql.NullString{String: fmt.Sprintf("%v", point), Valid: point != nil}, // Handle Point as NULL if not valid
// 		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
// 		UpdatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
// 	}

// 	// Save School entity in the database and get the UUID
// 	schoolUUID, err := service.registerRepository.SaveSchoolForRegister(school)
// 	if err != nil {
// 		return uuid.Nil, fmt.Errorf("failed to save school: %w", err)
// 	}

// 	return schoolUUID, nil
// }