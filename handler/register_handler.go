package handler

// import (
// 	"shuttle/errors"
// 	"shuttle/logger"
// 	"shuttle/models/dto"
// 	// "shuttle/models/entity"
// 	"shuttle/services"
// 	"shuttle/utils"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// )

// type RegisterHandlerInterface interface {
// 	AddUserRegister(c *fiber.Ctx) error
// 	ApproveUserRegister(c *fiber.Ctx) error
// }

// type registerHandler struct {
// 	registerService services.RegisterService
// 	schoolService services.SchoolService
// 	vehicleService services.VehicleService
// }

// func NewRegisterHttpHandler(registerService services.RegisterService, schoolService services.SchoolService, vehicleService services.VehicleService) RegisterHandlerInterface {
// 	return &registerHandler{
// 		registerService: registerService,
// 		schoolService: schoolService,
// 		vehicleService: vehicleService,
// 	}
// }

// func validateUserRoleDetailsForRegister(_ *fiber.Ctx, user *dto.UserRequestsDTO, handler registerHandler) error {
// 	switch user.Role {
// 	case dto.SuperAdmin:
// 		user.RoleCode = "SA"

// 	case dto.SchoolAdmin:
// 		details, err := parseDetails[dto.SchoolAdminDetailsRequestsDTO](user.Details)
// 		if err != nil {
// 			logger.LogError(err, "Invalid details format for SchoolAdmin", map[string]interface{}{
// 				"details": string(user.Details),
// 			})
// 			return errors.New("invalid details format for SchoolAdmin", 400)
// 		}

// 		if details.SchoolUUID == "" {
// 			// Jika SchoolUUID kosong, periksa apakah data sekolah baru tersedia
// 			if details.School.Name == "" || details.School.Address == "" {
// 				return errors.New("school is required for SchoolAdmin", 400)
// 			}
// 			// Jika data tersedia, lanjutkan proses tanpa error
// 		} else {
// 			// Jika SchoolUUID diisi, pastikan sekolah tersebut ada
// 			_, errSchool := handler.schoolService.GetSpecSchool(details.SchoolUUID)
// 			if errSchool != nil {
// 				return errors.New("school is not found", 404)
// 			}
// 		}

// 		user.RoleCode = "AS"

// 	case dto.Parent:
// 		if user.Details == nil {
// 			return errors.New("parent details are required", 400)
// 		}
// 		user.RoleCode = "P"

// 	case dto.Driver:
// 		details, err := parseDetails[dto.DriverDetailsRequestsDTO](user.Details)
// 		if err != nil {
// 			logger.LogError(err, "Invalid details format for Driver", map[string]interface{}{
// 				"details": string(user.Details),
// 			})
// 			return errors.New("invalid details format for Driver", 400)
// 		}

// 		if details.VehicleUUID != "" {
// 			_, errVehicle := handler.vehicleService.GetSpecVehicle(details.VehicleUUID)
// 			if errVehicle != nil {
// 				return errors.New("vehicle is not found", 404)
// 			}
// 		}

// 		if details.SchoolUUID != "" {
// 			_, errSchool := handler.schoolService.GetSpecSchool(details.SchoolUUID)
// 			if errSchool != nil {
// 				return errors.New("school is not found", 404)
// 			}
// 		}

// 		user.RoleCode = "D"

// 	default:
// 		return errors.New("invalid role specified", 400)
// 	}

// 	return nil
// }

// func (handler *registerHandler) AddUserRegister(c *fiber.Ctx) error {
// 	// Parsing request body ke DTO
// 	userReqDTO := new(dto.UserRequestsDTO)
// 	if err := c.BodyParser(userReqDTO); err != nil {
// 		return utils.BadRequestResponse(c, "Invalid request data", nil)
// 	}

// 	// Validasi struktur data request
// 	if err := utils.ValidateStruct(c, userReqDTO); err != nil {
// 		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
// 	}

// 	// Validasi role dan detail terkait
// 	if err := validateUserRoleDetailsForRegister(c, userReqDTO, *handler); err != nil {
// 		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
// 	}

// 	// Panggil service untuk menambahkan user
// 	if _, err := handler.registerService.AddUserRegister(*userReqDTO, ""); err != nil {
// 		if customErr, ok := err.(*errors.CustomError); ok {
// 			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
// 		}
// 		logger.LogError(err, "Failed to create user", nil)
// 		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
// 	}

// 	// Berikan respons sukses
// 	return utils.SuccessResponse(c, "User created successfully", nil)
// }

// func (handler *registerHandler) ApproveUserRegister(c *fiber.Ctx) error {
// 	userUUIDParam := c.Params("id")
// 	userUUID, err := uuid.Parse(userUUIDParam)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UUID tidak valid"})
// 	}

// 	username, ok := c.Locals("user_name").(string)
// 	if !ok || username == "" {
// 		// Handle missing or invalid username
// 		return utils.UnauthorizedResponse(c, "Unauthorized access: user_name not found", nil)
// 	}
// 	approverRole := c.Locals("role_code").(string)

// 	err = handler.registerService.ApproveUserRegister(userUUID, username, approverRole)
// 	if err != nil {
// 		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
// 	}

// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User berhasil diapprove"})
// }
