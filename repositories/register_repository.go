package repositories

// import (
// 	"encoding/json"
// 	"fmt"
// 	"shuttle/models/entity"

// 	"github.com/google/uuid"
// 	"github.com/jmoiron/sqlx"
// )

// type RegisterRepositoryInterface interface {
// 	BeginTransaction() (*sqlx.Tx, error)
// 	CheckEmailExistForRegister(uuid string, email string) (bool, error)
// 	CheckUsernameExistForRegister(uuid string, username string) (bool, error)
// 	GetUserByUUID(userUUID uuid.UUID) (entity.User, error)

// 	SaveUserRegister(tx *sqlx.Tx, userEntity entity.User) (uuid.UUID, error)
// 	SaveSchoolAdminRegisterDetails(tx *sqlx.Tx, details entity.SchoolAdminDetails, userUUID uuid.UUID, params interface{}) error
// 	SaveParentRegisterDetails(tx *sqlx.Tx, details entity.ParentDetails, userUUID uuid.UUID, params interface{}) error
// 	SaveDriverRegisterDetails(tx *sqlx.Tx, details entity.DriverDetails, userUUID uuid.UUID, params interface{}) error
// 	ApproveUserRegister(username string, userUUID uuid.UUID) error
// 	SaveSchoolForRegister(school entity.School) (uuid.UUID, error)
// }

// type registerRepository struct {
// 	DB *sqlx.DB
// }

// func NewRegisterRepository(DB *sqlx.DB) RegisterRepositoryInterface {
// 	return &registerRepository{
// 		DB: DB,
// 	}
// }

// func (r *registerRepository) BeginTransaction() (*sqlx.Tx, error) {
// 	tx, err := r.DB.Beginx()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return tx, nil
// }

// func (r *registerRepository) CheckEmailExistForRegister(uuid string, email string) (bool, error) {
// 	var count int
// 	query := `SELECT COUNT(user_id) FROM users WHERE user_email = $1 AND deleted_at IS NULL`

// 	if uuid != "" {
// 		query += ` AND user_uuid != $2`
// 		if err := r.DB.Get(&count, query, email, uuid); err != nil {
// 			return false, err
// 		}
// 	} else {
// 		if err := r.DB.Get(&count, query, email); err != nil {
// 			return false, err
// 		}
// 	}

// 	return count > 0, nil
// }

// func (r *registerRepository) CheckUsernameExistForRegister(uuid string, username string) (bool, error) {
// 	var count int
// 	query := `SELECT COUNT(user_id) FROM users WHERE user_username = $1 AND deleted_at IS NULL`

// 	if uuid != "" {
// 		query += ` AND user_uuid != $2`
// 		if err := r.DB.Get(&count, query, username, uuid); err != nil {
// 			return false, err
// 		}
// 	} else {
// 		if err := r.DB.Get(&count, query, username); err != nil {
// 			return false, err
// 		}
// 	}

// 	return count > 0, nil
// }

// func (r *registerRepository) GetUserByUUID(userUUID uuid.UUID) (entity.User, error) {
//     var user entity.User
//     query := `
//         SELECT user_uuid, user_username, user_email, user_role, user_role_code, user_register_status 
//         FROM users 
//         WHERE user_uuid = $1
//     `
//     err := r.DB.Get(&user, query, userUUID)
//     if err != nil {
//         return entity.User{}, err
//     }
//     return user, nil
// }


// func (r *registerRepository) SaveUserRegister(tx *sqlx.Tx, userEntity entity.User) (uuid.UUID, error) {
// 	query := `
// 		INSERT INTO users (user_id, user_uuid, user_username, user_email, user_password, user_role, user_role_code, user_register_status)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
// 		RETURNING user_uuid`
// 	var userUUID uuid.UUID
// 	err := tx.QueryRow(query, userEntity.ID, userEntity.UUID, userEntity.Username, userEntity.Email, userEntity.Password, userEntity.Role, userEntity.RoleCode, userEntity.RegisterStatus).Scan(&userUUID)
// 	if err != nil {
// 		return uuid.Nil, err
// 	}
// 	return userUUID, nil
// }

// func (r *registerRepository) SaveSchoolAdminRegisterDetails(tx *sqlx.Tx, details entity.SchoolAdminDetails, userUUID uuid.UUID, params interface{}) error {
// 	details.UserUUID = userUUID
// 	query := `
//         INSERT INTO school_admin_details 
//         (user_uuid, school_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address) 
//         VALUES (:user_uuid, :school_uuid, :user_picture, :user_first_name, :user_last_name, :user_gender, :user_phone, :user_address)
//     `
// 	params = details
// 	_, err := tx.NamedExec(query, params)
// 	return err
// }

// func (r *registerRepository) SaveParentRegisterDetails(tx *sqlx.Tx, details entity.ParentDetails, userUUID uuid.UUID, params interface{}) error {
// 	details.UserUUID = userUUID
// 	query := `
//         INSERT INTO parent_details 
//         (user_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address) 
//         VALUES (:user_uuid, :user_picture, :user_first_name, :user_last_name, :user_gender, :user_phone, :user_address)
//     `
// 	params = details
// 	_, err := tx.NamedExec(query, params)
// 	return err
// }

// func (r *registerRepository) SaveDriverRegisterDetails(tx *sqlx.Tx, details entity.DriverDetails, userUUID uuid.UUID, params interface{}) error {
// 	details.UserUUID = userUUID

// 	if details.SchoolUUID == nil || *details.SchoolUUID == uuid.Nil {
// 		details.SchoolUUID = nil
// 	}
// 	if details.VehicleUUID == nil || *details.VehicleUUID == uuid.Nil {
// 		details.VehicleUUID = nil
// 	}

// 	query := `
// 		INSERT INTO driver_details 
// 		(user_uuid, school_uuid, vehicle_uuid, user_picture, user_first_name, user_last_name, user_gender, user_phone, user_address, user_license_number) 
// 		VALUES (:user_uuid, :school_uuid, :vehicle_uuid, :user_picture, :user_first_name, :user_last_name, :user_gender, :user_phone, :user_address, :user_license_number)
// 	`
// 	params = details
// 	_, err := tx.NamedExec(query, params)
// 	if err != nil {
// 		return err
// 	}

// 	if details.VehicleUUID != nil {
// 		return r.UpdateDriverUUIDOnVehicles(tx, userUUID, *details.VehicleUUID)
// 	}
// 	return nil
// }

// func (r *registerRepository) UpdateDriverUUIDOnVehicles(tx *sqlx.Tx, userUUID uuid.UUID, vehicleUUID uuid.UUID) error {
// 	var userUUIDParam interface{}
// 	if userUUID == uuid.Nil {
// 		userUUIDParam = nil
// 	} else {
// 		userUUIDParam = userUUID
// 	}

// 	query := `
//         UPDATE vehicles
//         SET driver_uuid = $1
//         WHERE vehicle_uuid = $2
// 		`
// 	_, err := tx.Exec(query, userUUIDParam, vehicleUUID)
// 	return err
// }

// func (r *registerRepository) ApproveUserRegister(username string, userUUID uuid.UUID) error {
//     query := `
//         UPDATE users 
//         SET user_register_status = 'APPROVED', 
//             approved_by = $1, 
//             approved_at = NOW() 
//         WHERE user_uuid = $2 AND user_register_status = 'PENDING'
//     `
//     _, err := r.DB.Exec(query, username, userUUID)
//     return err
// }

// func (r *registerRepository) SaveSchoolForRegister(school entity.School) (uuid.UUID, error) {
// 	// Mengonversi Point menjadi JSON string jika valid
// 	var point string
// 	if school.Point.Valid {
// 		// Jika Point ada, kita marshal ke JSON string
// 		pointBytes, err := json.Marshal(school.Point.String)
// 		if err != nil {
// 			return uuid.Nil, fmt.Errorf("gagal mengonversi point ke JSON: %w", err)
// 		}
// 		point = string(pointBytes) // Menyimpan sebagai string JSON
// 	} else {
// 		point = "" // NULL jika tidak valid
// 	}

// 	// Query untuk menyimpan data ke dalam tabel schools
// 	query := `
// 		INSERT INTO schools (school_id, school_uuid, school_name, school_address, school_contact, school_email, school_description, school_point, created_by)
// 		VALUES (:school_id, :school_uuid, :school_name, :school_address, :school_contact, :school_email, :school_description, :school_point, :created_by)
// 		RETURNING school_uuid`

// 	// Menyusun parameter dengan map[string]interface{}
// 	params := map[string]interface{}{
// 		"school_id":        school.ID,
// 		"school_uuid":      school.UUID,
// 		"school_name":      school.Name,
// 		"school_address":   school.Address,
// 		"school_contact":   school.Contact,
// 		"school_email":     school.Email,
// 		"school_description": school.Description,
// 		"school_point":     point, // Menyimpan JSON string
// 		"created_by":       school.CreatedBy,
// 	}

// 	// Menggunakan DB Get untuk eksekusi query dan mendapatkan UUID sekolah yang baru
// 	var newSchoolUUID uuid.UUID
// 	err := r.DB.Get(&newSchoolUUID, query, params)
// 	if err != nil {
// 		return uuid.Nil, fmt.Errorf("gagal menyimpan sekolah: %w", err)
// 	}

// 	// Mengembalikan UUID dari sekolah yang baru saja disimpan
// 	return newSchoolUUID, nil
// }