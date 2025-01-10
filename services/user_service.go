package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"shuttle/errors"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserServiceInterface interface {
	////////////////////////////////////// TEMPORARY //////////////////////////////////////
	GetDriverDetailsByUUID(driverUUID uuid.UUID) (entity.DriverDetails, error)
	GetSchoolAdminDetailsByUUID(schoolAdminUUID uuid.UUID) (entity.SchoolAdminDetails, error)
	GetSpecDriverForPermittedSchool(driverUUID string, schoolUUID string) (dto.UserResponseDTO, error)
	////////////////////////////////////// TEMPORARY //////////////////////////////////////

	GetAllSuperAdmin(page int, limit int, sortField string, sortDirection string) ([]dto.UserResponseDTO, int, error)
	GetSpecSuperAdmin(uuid string) (dto.UserResponseDTO, error)
	GetAllSchoolAdmin(page int, limit int, sortField string, sortDirection string) ([]dto.UserResponseDTO, int, error)
	GetSpecSchoolAdmin(uuid string) (dto.UserResponseDTO, error)

	GetAllDriverFromAllSchools(page int, limit int, sortField string, sortDirection string) ([]dto.UserResponseDTO, int, error)
	GetAllDriverForPermittedSchool(page int, limit int, sortField string, sortDirection string, schoolUUID string) ([]dto.UserResponseDTO, int, error)
	GetSpecDriverFromAllSchools(uuid string) (dto.UserResponseDTO, error)

	AddUser(req dto.UserRequestsDTO, user_name string) (uuid.UUID, error)
	UpdateUser(id string, user dto.UserRequestsDTO, user_name string, file []byte) error

	DeleteSuperAdmin(id string, user_name string) error
	DeleteSchoolAdmin(id string, user_name string) error
	DeleteDriver(id string, user_name string) error

	// GetSpecUser(id string) (entity.User, error)
	GetSpecUserWithDetails(id string) (UserWithDetails, error)

	CheckPermittedSchoolAccess(userUUID string) (string, error)
}

type UserService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewUserService(userRepository repositories.UserRepositoryInterface) UserService {
	return UserService{
		userRepository: userRepository,
	}
}

func (service *UserService) GetDriverDetailsByUUID(driverUUID uuid.UUID) (entity.DriverDetails, error) {
	return service.userRepository.FetchDriverDetails(driverUUID)
}

func (service *UserService) GetSchoolAdminDetailsByUUID(schoolAdminUUID uuid.UUID) (entity.SchoolAdminDetails, error) {
	return service.userRepository.FetchSchoolAdminDetails(schoolAdminUUID)
}

func (service *UserService) GetAllSuperAdmin(page int, limit int, sortField, sortDirection string) ([]dto.UserResponseDTO, int, error) {
	offset := (page - 1) * limit

	users, err := service.userRepository.FetchAllSuperAdmins(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.userRepository.CountSuperAdmin()
	if err != nil {
		return nil, 0, err
	}

	var usersDTO []dto.UserResponseDTO

	for _, user := range users {
		userDTO := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}

		superAdminDetails, err := service.userRepository.FetchSuperAdminDetails(user.UUID)
		if err != nil {
			return nil, 0, err
		}

		detailsJSON, err := json.Marshal(dto.SuperAdminDetailsResponseDTO{
			FirstName: superAdminDetails.FirstName,
			LastName:  superAdminDetails.LastName,
			Gender:    dto.Gender(superAdminDetails.Gender),
			Phone:     superAdminDetails.Phone,
		})
		if err != nil {
			return nil, 0, err
		}
		userDTO.Details = detailsJSON

		usersDTO = append(usersDTO, userDTO)
	}

	return usersDTO, total, nil
}

func (service *UserService) GetAllSchoolAdmin(page int, limit int, sortField, sortDirection string) ([]dto.UserResponseDTO, int, error) {
	offset := (page - 1) * limit

	users, school, err := service.userRepository.FetchAllSchoolAdmins(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.userRepository.CountSchoolAdmin()
	if err != nil {
		return nil, 0, err
	}

	var usersDTO []dto.UserResponseDTO

	for _, user := range users {
		userDTO := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}

		schoolAdminDetails, err := service.userRepository.FetchSchoolAdminDetails(user.UUID)
		if err != nil {
			return nil, 0, err
		}

		detailsJSON, err := json.Marshal(dto.SchoolAdminDetailsResponseDTO{
			SchoolName: school.Name,
			Picture:    schoolAdminDetails.Picture,
			FirstName:  schoolAdminDetails.FirstName,
			LastName:   schoolAdminDetails.LastName,
			Gender:     dto.Gender(schoolAdminDetails.Gender),
			Phone:      schoolAdminDetails.Phone,
		})
		if err != nil {
			return nil, 0, err
		}
		userDTO.Details = detailsJSON

		usersDTO = append(usersDTO, userDTO)
	}

	return usersDTO, total, nil
}

func (service *UserService) GetAllDriverFromAllSchools(page int, limit int, sortField string, sortDirection string) ([]dto.UserResponseDTO, int, error) {
	offset := (page - 1) * limit
	log.Println("Fetching all drivers with pagination parameters - page:", page, "limit:", limit, "sortField:", sortField, "sortDirection:", sortDirection)

	users, school, vehicle, err := service.userRepository.FetchAllDrivers(offset, limit, sortField, sortDirection)
	if err != nil {
		log.Println("Error fetching drivers:", err)
		return nil, 0, err
	}

	log.Println("Successfully fetched drivers. Total drivers:", len(users))

	total, err := service.userRepository.CountAllPermittedDriver("")
	if err != nil {
		log.Println("Error counting permitted drivers:", err)
		return nil, 0, err
	}
	log.Println("Total permitted drivers:", total)

	var usersDTO []dto.UserResponseDTO

	for _, user := range users {
		log.Println("Processing user:", user.UUID.String())

		userDTO := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}

		driverDetails, err := service.userRepository.FetchDriverDetails(user.UUID)
		if err != nil {
			log.Println("Error fetching driver details for user:", user.UUID.String(), err)
			return nil, 0, err
		}

		log.Println("Successfully fetched driver details for user:", user.UUID.String())

		var vehicleDetails string
		if vehicle.VehicleNumber == "N/A" || vehicle.UUID == uuid.Nil {
			vehicleDetails = "N/A"
		} else {
			vehicleDetails = fmt.Sprintf("%s (%s)", vehicle.VehicleNumber, vehicle.VehicleName)
		}

		log.Println("Formatted vehicle details:", vehicleDetails)

		detailsJSON, err := json.Marshal(dto.DriverDetailsResponseDTO{
			SchoolName:    school.Name,
			VehicleNumber: vehicleDetails,
			Picture:       driverDetails.Picture,
			FirstName:     driverDetails.FirstName,
			LastName:      driverDetails.LastName,
			Gender:        dto.Gender(driverDetails.Gender),
			Phone:         driverDetails.Phone,
			Address:       driverDetails.Address,
			LicenseNumber: driverDetails.LicenseNumber,
		})
		if err != nil {
			log.Println("Error marshalling driver details to JSON:", err)
			return nil, 0, err
		}

		log.Println("Successfully marshalled driver details to JSON")

		userDTO.Details = detailsJSON
		usersDTO = append(usersDTO, userDTO)

		log.Println("User DTO created for user:", user.UUID.String())
	}

	log.Println("Returning total users and user DTOs.")
	return usersDTO, total, nil
}

func (service *UserService) GetAllDriverForPermittedSchool(page int, limit int, sortField string, sortDirection string, schoolUUID string) ([]dto.UserResponseDTO, int, error) {
	offset := (page - 1) * limit
	log.Println("Fetching all drivers for school with UUID:", schoolUUID, "page:", page, "limit:", limit, "sortField:", sortField, "sortDirection:", sortDirection)

	users, school, vehicle, err := service.userRepository.FetchAllDriversForPermittedSchool(offset, limit, sortField, sortDirection, schoolUUID)
	if err != nil {
		log.Println("Error fetching drivers for school:", err)
		return nil, 0, err
	}

	log.Println("Successfully fetched drivers for school:", schoolUUID, "Total drivers:", len(users))

	total, err := service.userRepository.CountAllPermittedDriver(schoolUUID)
	if err != nil {
		log.Println("Error counting permitted drivers for school:", err)
		return nil, 0, err
	}
	log.Println("Total permitted drivers for school:", schoolUUID, "Total:", total)

	var usersDTO []dto.UserResponseDTO

	for _, user := range users {
		log.Println("Processing user:", user.UUID.String(), "for school:", schoolUUID)

		userDTO := dto.UserResponseDTO{
			UUID:       user.UUID.String(),
			Username:   user.Username,
			Email:      user.Email,
			Status:     user.Status,
			LastActive: safeTimeFormat(user.LastActive),
			CreatedAt:  safeTimeFormat(user.CreatedAt),
		}

		driverDetails, err := service.userRepository.FetchDriverDetails(user.UUID)
		if err != nil {
			log.Println("Error fetching driver details for user:", user.UUID.String(), err)
			return nil, 0, err
		}

		log.Println("Successfully fetched driver details for user:", user.UUID.String())

		detailsJSON, err := json.Marshal(dto.DriverDetailsResponseDTO{
			SchoolName:    school.Name,
			VehicleNumber: vehicle.VehicleNumber,
			Picture:       driverDetails.Picture,
			FirstName:     driverDetails.FirstName,
			LastName:      driverDetails.LastName,
			Gender:        dto.Gender(driverDetails.Gender),
			Phone:         driverDetails.Phone,
			Address:       driverDetails.Address,
			LicenseNumber: driverDetails.LicenseNumber,
		})
		if err != nil {
			log.Println("Error marshalling driver details to JSON:", err)
			return nil, 0, err
		}

		log.Println("Successfully marshalled driver details to JSON")

		userDTO.Details = detailsJSON
		usersDTO = append(usersDTO, userDTO)

		log.Println("User DTO created for user:", user.UUID.String())
	}

	log.Println("Returning total users and user DTOs for permitted school.")
	return usersDTO, total, nil
}


func (service *UserService) GetSpecSuperAdmin(uuid string) (dto.UserResponseDTO, error) {
	user, err := service.userRepository.FetchSpecSuperAdmin(uuid)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO := dto.UserResponseDTO{
		UUID:       user.UUID.String(),
		Username:   user.Username,
		Email:      user.Email,
		Status:     user.Status,
		LastActive: safeTimeFormat(user.LastActive),
		CreatedAt:  safeTimeFormat(user.CreatedAt),
		CreatedBy:  safeStringFormat(user.CreatedBy),
		UpdatedAt:  safeTimeFormat(user.UpdatedAt),
		UpdatedBy:  safeStringFormat(user.UpdatedBy),
	}

	superAdminDetails, err := service.userRepository.FetchSuperAdminDetails(user.UUID)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	detailsJSON, err := json.Marshal(dto.SuperAdminDetailsResponseDTO{
		Picture:   superAdminDetails.Picture,
		FirstName: superAdminDetails.FirstName,
		LastName:  superAdminDetails.LastName,
		Gender:    dto.Gender(superAdminDetails.Gender),
		Phone:     superAdminDetails.Phone,
		Address:   superAdminDetails.Address,
	})
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO.Details = detailsJSON

	return userDTO, nil
}

func (service *UserService) GetSpecSchoolAdmin(id string) (dto.UserResponseDTO, error) {
	user, school, err := service.userRepository.FetchSpecSchoolAdmin(id)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO := dto.UserResponseDTO{
		UUID:       user.UUID.String(),
		Username:   user.Username,
		Email:      user.Email,
		Status:     user.Status,
		LastActive: safeTimeFormat(user.LastActive),
		CreatedAt:  safeTimeFormat(user.CreatedAt),
		CreatedBy:  safeStringFormat(user.CreatedBy),
		UpdatedAt:  safeTimeFormat(user.UpdatedAt),
		UpdatedBy:  safeStringFormat(user.UpdatedBy),
	}

	schoolAdminDetails, err := service.userRepository.FetchSchoolAdminDetails(user.UUID)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	detailsJSON, err := json.Marshal(dto.SchoolAdminDetailsResponseDTO{
		SchoolUUID: schoolAdminDetails.SchoolUUID.String(),
		SchoolName: school.Name,
		Picture:    schoolAdminDetails.Picture,
		FirstName:  schoolAdminDetails.FirstName,
		LastName:   schoolAdminDetails.LastName,
		Gender:     dto.Gender(schoolAdminDetails.Gender),
		Phone:      schoolAdminDetails.Phone,
		Address:    schoolAdminDetails.Address,
	})
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO.Details = detailsJSON

	return userDTO, nil
}

func (service *UserService) GetSpecDriverFromAllSchools(id string) (dto.UserResponseDTO, error) {
	user, school, vehicle, err := service.userRepository.FetchSpecDriverFromAllSchools(id)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO := dto.UserResponseDTO{
		UUID:       user.UUID.String(),
		Username:   user.Username,
		Email:      user.Email,
		Status:     user.Status,
		LastActive: safeTimeFormat(user.LastActive),
		CreatedAt:  safeTimeFormat(user.CreatedAt),
		CreatedBy:  safeStringFormat(user.CreatedBy),
		UpdatedAt:  safeTimeFormat(user.UpdatedAt),
		UpdatedBy:  safeStringFormat(user.UpdatedBy),
	}

	driverDetails, err := service.userRepository.FetchDriverDetails(user.UUID)
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	var vehicleDetails, schoolUUID, vehicleUUID string
	if vehicle.VehicleNumber == "N/A" || vehicle.UUID == uuid.Nil {
		vehicleDetails = "N/A"
	} else {
		vehicleDetails = fmt.Sprintf("%s (%s)", vehicle.VehicleNumber, vehicle.VehicleName)
	}

	if school.UUID == uuid.Nil {
		schoolUUID = "N/A"
	} else {
		schoolUUID = school.UUID.String()
	}

	if vehicle.UUID == uuid.Nil {
		vehicleUUID = "N/A"
	} else {
		vehicleUUID = vehicle.UUID.String()
	}

	detailsJSON, err := json.Marshal(dto.DriverDetailsResponseDTO{
		SchoolUUID:    schoolUUID,
		SchoolName:    school.Name,
		VehicleUUID:   vehicleUUID,
		VehicleNumber: vehicleDetails,
		Picture:       driverDetails.Picture,
		FirstName:     driverDetails.FirstName,
		LastName:      driverDetails.LastName,
		Gender:        dto.Gender(driverDetails.Gender),
		Phone:         driverDetails.Phone,
		Address:       driverDetails.Address,
		LicenseNumber: driverDetails.LicenseNumber,
	})
	if err != nil {
		return dto.UserResponseDTO{}, err
	}

	userDTO.Details = detailsJSON

	return userDTO, nil
}

func (service *UserService) GetSpecDriverForPermittedSchool(driverUUID string, schoolUUID string) (dto.UserResponseDTO, error) {
	log.Println("Fetching spec driver for school...")
	log.Printf("Driver UUID: %s, School UUID: %s", driverUUID, schoolUUID)
	
	user, school, vehicle, err := service.userRepository.FetchSpecDriverForPermittedSchool(driverUUID, schoolUUID)
	if err != nil {
		log.Println("Error fetching driver data:", err)
		return dto.UserResponseDTO{}, err
	}
	log.Println("Fetched user, school, and vehicle data successfully.")

	log.Println("User data fetched:")
	log.Printf("User UUID: %s, Username: %s, Email: %s", user.UUID, user.Username, user.Email)

	userDTO := dto.UserResponseDTO{
		UUID:       user.UUID.String(),
		Username:   user.Username,
		Email:      user.Email,
		Status:     user.Status,
		LastActive: safeTimeFormat(user.LastActive),
		CreatedAt:  safeTimeFormat(user.CreatedAt),
		CreatedBy:  safeStringFormat(user.CreatedBy),
		UpdatedAt:  safeTimeFormat(user.UpdatedAt),
		UpdatedBy:  safeStringFormat(user.UpdatedBy),
	}

	log.Println("Driver details fetching...")
	driverDetails, err := service.userRepository.FetchDriverDetails(user.UUID)
	if err != nil {
		log.Println("Error fetching driver details:", err)
		return dto.UserResponseDTO{}, err
	}
	log.Printf("Fetched driver details: First Name: %s, Last Name: %s, Phone: %s", driverDetails.FirstName, driverDetails.LastName, driverDetails.Phone)

	var vehicleDetails, vehicleUUID string
	if vehicle.VehicleNumber == "N/A" || vehicle.UUID == uuid.Nil {
		vehicleDetails = "N/A"
	} else {
		vehicleDetails = fmt.Sprintf("%s (%s)", vehicle.VehicleNumber, vehicle.VehicleName)
	}
	log.Printf("Vehicle details: %s", vehicleDetails)

	if vehicle.UUID == uuid.Nil {
		vehicleUUID = "N/A"
	} else {
		vehicleUUID = vehicle.UUID.String()
	}
	log.Printf("Vehicle UUID: %s", vehicleUUID)

	// Marshal details into JSON
	log.Println("Marshaling driver details into JSON...")
	driverDetailsJSON, err := json.Marshal(dto.DriverDetailsResponseDTO{
		SchoolUUID:    schoolUUID,
		SchoolName:    school.Name,
		VehicleUUID:   vehicleUUID,
		VehicleNumber: vehicleDetails,
		Picture:       driverDetails.Picture,
		FirstName:     driverDetails.FirstName,
		LastName:      driverDetails.LastName,
		Gender:        dto.Gender(driverDetails.Gender),
		Phone:         driverDetails.Phone,
		Address:       driverDetails.Address,
		LicenseNumber: driverDetails.LicenseNumber,
	})
	if err != nil {
		log.Println("Error marshaling driver details:", err)
		return dto.UserResponseDTO{}, err
	}

	// Log the JSON in readable format
	log.Println("Pretty-printing driver details JSON...")
	prettyJSON, err := json.MarshalIndent(driverDetailsJSON, "", "  ")
	if err != nil {
		log.Println("Error pretty-printing driver details JSON:", err)
		return dto.UserResponseDTO{}, err
	}
	log.Printf("Driver Details JSON: %s", string(prettyJSON))

	userDTO.Details = driverDetailsJSON
	log.Println("Driver details added to user DTO.")

	// Log the final DTO
	log.Printf("User details: %+v", userDTO)

	return userDTO, nil
}

func (s *UserService) AddUser(req dto.UserRequestsDTO, user_name string) (uuid.UUID, error) {
	exists, err := s.userRepository.CheckEmailExist("", req.Email)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, errors.New("email already exists", 409)
	}

	exists, err = s.userRepository.CheckUsernameExist("", req.Username)
	if err != nil {
		return uuid.Nil, err
	}
	if exists {
		return uuid.Nil, errors.New("username already exists", 409)
	}

	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			return uuid.Nil, err
		}
		req.Password = hashedPassword
	}

	tx, err := s.userRepository.BeginTransaction()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	var transactionErr error
	defer func() {
		if transactionErr != nil {
			tx.Rollback()
		} else {
			transactionErr = tx.Commit()
		}
	}()

	userEntity := entity.User{
		ID:        time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		UUID:      uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		Role:      entity.Role(req.Role),
		RoleCode:  req.RoleCode,
		CreatedBy: sql.NullString{String: user_name, Valid: user_name != ""},
	}

	userUUID, err := s.userRepository.SaveUser(tx, userEntity)
	if err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return uuid.Nil, errors.New(customErr.Message, customErr.StatusCode)
		}
		transactionErr = fmt.Errorf("error saving user: %w", err)
		return uuid.Nil, transactionErr
	}

	if err := s.saveRoleDetails(tx, userEntity.UUID, req); err != nil {
		transactionErr = fmt.Errorf("error saving role details: %w", err)
		return uuid.Nil, transactionErr
	}

	return userUUID, nil
}

func (s *UserService) UpdateUser(id string, req dto.UserRequestsDTO, username string, file []byte) error {
	tx, err := s.userRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	exists, err := s.userRepository.CheckEmailExist(id, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists", 409)
	}

	exists, err = s.userRepository.CheckUsernameExist(id, req.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists", 409)
	}

	userData := entity.User{
		Username:  req.Username,
		Email:     req.Email,
		Role:      entity.Role(req.Role),
		RoleCode:  req.RoleCode,
		UpdatedBy: sql.NullString{String: username, Valid: username != ""},
	}

	if err := s.userRepository.UpdateUser(tx, userData, id); err != nil {
		return err
	}

	if err := s.updateRoleDetails(tx, req, id); err != nil {
		logger.LogError(err, "error updating role details", map[string]interface{}{})
		return fmt.Errorf("error updating role details: %w", err)
	}

	return nil
}

func (service *UserService) DeleteSuperAdmin(id string, user_name string) error {
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("user not found", 404)
	}

	err = service.userRepository.DeleteSuperAdmin(tx, parsedUUID, user_name)
	if err != nil {
		return errors.New("user not found", 404)
	}

	return nil
}

func (service *UserService) DeleteSchoolAdmin(id string, user_name string) error {
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	var transactionErr error
	defer func() {
		if transactionErr != nil {
			tx.Rollback()
		} else {
			transactionErr = tx.Commit()
		}
	}()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		transactionErr = errors.New("user not found", 404)
		return transactionErr
	}

	err = service.userRepository.DeleteSchoolAdmin(tx, parsedUUID, user_name)
	if err != nil {
		transactionErr = errors.New("user not found", 404)
		return transactionErr
	}

	return nil
}

func (service *UserService) DeleteDriver(id string, user_name string) error {
	tx, err := service.userRepository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	var transactionErr error
	defer func() {
		if transactionErr != nil {
			tx.Rollback()
		} else {
			transactionErr = tx.Commit()
		}
	}()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		transactionErr = errors.New("user not found", 404)
		return transactionErr
	}

	err = service.userRepository.DeleteDriver(tx, parsedUUID, user_name)
	if err != nil {
		transactionErr = errors.New("user not found", 404)
		return transactionErr
	}

	return nil
}

type UserWithDetails struct {
	User               entity.User                `json:"user"`
	SuperAdminDetails  *entity.SuperAdminDetails  `json:"super_admin_details,omitempty"`
	SchoolAdminDetails *entity.SchoolAdminDetails `json:"school_admin_details,omitempty"`
	DriverDetails      *entity.DriverDetails      `json:"driver_details,omitempty"`
	ParentDetails      *entity.ParentDetails      `json:"parent_details,omitempty"`
}

func (service *UserService) GetSpecUserWithDetails(id string) (UserWithDetails, error) {
	user, err := service.userRepository.FetchSpecificUser(id)
	if err != nil {
		return UserWithDetails{}, err
	}

	userWithDetails := UserWithDetails{
		User: user,
	}

	switch user.RoleCode {
	case "SA":
		superAdminDetails, err := service.userRepository.FetchSuperAdminDetails(user.UUID)
		if err != nil {
			return UserWithDetails{}, err
		}
		userWithDetails.SuperAdminDetails = &superAdminDetails

	case "AS":
		schoolAdminDetails, err := service.userRepository.FetchSchoolAdminDetails(user.UUID)
		if err != nil {
			return UserWithDetails{}, err
		}
		userWithDetails.SchoolAdminDetails = &schoolAdminDetails

	case "P":
		parentDetails, err := service.userRepository.FetchParentDetails(user.UUID)
		if err != nil {
			return UserWithDetails{}, err
		}
		userWithDetails.ParentDetails = &parentDetails

	case "D":
		driverDetails, err := service.userRepository.FetchDriverDetails(user.UUID)
		if err != nil {
			return UserWithDetails{}, err
		}
		userWithDetails.DriverDetails = &driverDetails

	default:
		return UserWithDetails{}, errors.New("invalid role code", 0)
	}

	return userWithDetails, nil
}

// func (service *UserService) GetSpecUser(id string) (entity.User, error) {
// 	db, err := databases.PostgresConnection()
// 	if err != nil {
// 		return entity.User{}, err
// 	}

// 	idInt, err := strconv.ParseInt(id, 10, 64)
// 	if err != nil {
// 		return entity.User{}, errors.New("invalid user id", 0)
// 	}

// 	var user entity.User
// 	query := `
// 		SELECT * FROM users WHERE user_id = $1
// 	`

// 	err = db.Get(&user, query, idInt)
// 	if err != nil {
// 		return entity.User{}, err
// 	}

// 	return user, nil
// }

func (service *UserService) CheckPermittedSchoolAccess(userUUID string) (string, error) {
	schoolUUID, err := service.userRepository.FetchPermittedSchoolAccess(userUUID)
	if err != nil {
		return "", err
	}

	return schoolUUID, nil
}

func (s *UserService) saveRoleDetails(tx *sqlx.Tx, userUUID uuid.UUID, req dto.UserRequestsDTO) error {
	switch entity.Role(req.Role) {
	case entity.SuperAdmin:
		details := entity.SuperAdminDetails{
			Picture:   req.Picture,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    entity.Gender(req.Gender),
			Phone:     req.Phone,
			Address:   req.Address,
		}
		return s.userRepository.SaveSuperAdminDetails(tx, details, userUUID, nil)

	case entity.SchoolAdmin:
		parsedDetails, err := parseDetails[dto.SchoolAdminDetailsRequestsDTO](req.Details)
		if err != nil {
			return errors.New("invalid school admin details format: "+err.Error(), 400)
		}
		schoolDetails := entity.SchoolAdminDetails{
			SchoolUUID: uuid.MustParse(parsedDetails.SchoolUUID),
			Picture:    req.Picture,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Gender:     entity.Gender(req.Gender),
			Phone:      req.Phone,
			Address:    req.Address,
		}
		return s.userRepository.SaveSchoolAdminDetails(tx, schoolDetails, userUUID, nil)

	case entity.Parent:
		details := entity.ParentDetails{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    entity.Gender(req.Gender),
			Phone:     req.Phone,
			Address:   req.Address,
		}
		return s.userRepository.SaveParentDetails(tx, details, userUUID, nil)

	case entity.Driver:
		parsedDetails, err := parseDetails[dto.DriverDetailsRequestsDTO](req.Details)
		if err != nil {
			return errors.New("invalid driver details format: "+err.Error(), 400)
		}
		
		fmt.Printf("driver details: %+v\n", parsedDetails)
		println("vehicleUUID: ", parsedDetails.VehicleUUID)
		driverDetails := entity.DriverDetails{
			SchoolUUID:    parseSafeUUID(parsedDetails.SchoolUUID),
			VehicleUUID:   parseSafeUUID(parsedDetails.VehicleUUID),
			Picture:       req.Picture,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			Gender:        entity.Gender(req.Gender),
			Phone:         req.Phone,
			Address:       req.Address,
			LicenseNumber: parsedDetails.LicenseNumber,
		}
		return s.userRepository.SaveDriverDetails(tx, driverDetails, userUUID, nil)

	default:
		return errors.New("invalid role", 400)
	}
}

func (s *UserService) updateRoleDetails(tx *sqlx.Tx, req dto.UserRequestsDTO, id string) error {
	switch entity.Role(req.Role) {
	case entity.SuperAdmin:
		details := entity.SuperAdminDetails{
			Picture:   req.Picture,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    entity.Gender(req.Gender),
			Phone:     req.Phone,
			Address:   req.Address,
		}
		if err := s.userRepository.UpdateSuperAdminDetails(tx, details, id); err != nil {
			return err
		}

	case entity.SchoolAdmin:
		var details dto.SchoolAdminDetailsRequestsDTO
		if err := json.Unmarshal(req.Details, &details); err != nil {
			return errors.New("invalid school admin details format", 400)
		}

		schoolDetails := entity.SchoolAdminDetails{
			SchoolUUID: uuid.MustParse(details.SchoolUUID),
			Picture:    req.Picture,
			FirstName:  req.FirstName,
			LastName:   req.LastName,
			Gender:     entity.Gender(req.Gender),
			Phone:      req.Phone,
			Address:    req.Address,
		}
		if err := s.userRepository.UpdateSchoolAdminDetails(tx, schoolDetails, id); err != nil {
			return err
		}

	case entity.Parent:
		details := entity.ParentDetails{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Gender:    entity.Gender(req.Gender),
			Phone:     req.Phone,
			Address:   req.Address,
		}
		if err := s.userRepository.UpdateParentDetails(tx, details, id); err != nil {
			return err
		}			

	case entity.Driver:
		var details dto.DriverDetailsRequestsDTO
		if err := json.Unmarshal(req.Details, &details); err != nil {
			return errors.New("invalid driver details format", 400)
		}

		driverDetails := entity.DriverDetails{
			SchoolUUID:    parseSafeUUID(details.SchoolUUID),
			VehicleUUID:   parseSafeUUID(details.VehicleUUID),
			Picture:       req.Picture,
			FirstName:     req.FirstName,
			LastName:      req.LastName,
			Gender:        entity.Gender(req.Gender),
			Phone:         req.Phone,
			Address:       req.Address,
			LicenseNumber: details.LicenseNumber,
		}
		parsedUUID, err := uuid.Parse(id)
		if err != nil {
			return fmt.Errorf("invalid UUID: %w", err)
		}
		if err := s.userRepository.UpdateDriverDetails(tx, driverDetails, parsedUUID); err != nil {
			return err
		}

	default:
		return errors.New("invalid role", 400)
	}

	return nil
}

func parseSafeUUID(id string) *uuid.UUID {
	if id == "" || id == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return &parsedUUID
}

func parseDetails[T any](details json.RawMessage) (T, error) {
	var parsedDetails T

	err := json.Unmarshal(details, &parsedDetails)
	if err != nil {
		return parsedDetails, fmt.Errorf("failed to unmarshal details to struct: %w", err)
	}

	return parsedDetails, nil
}