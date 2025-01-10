package handler

import (
	"fmt"
	"log"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandlerInterface interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	GetMyProfile(c *fiber.Ctx) error
	IssueNewAccessToken(c *fiber.Ctx) error
	AddDeviceToken(c *fiber.Ctx) error
}

type authHandler struct {
	authService services.AuthService
}

func NewAuthHttpHandler(authService services.AuthService) AuthHandlerInterface {
	return &authHandler{
		authService: authService,
	}
}

func (handler *authHandler) Login(c *fiber.Ctx) error {
	loginRequest := new(dto.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	userDataOnLogin, err := handler.authService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		logger.LogError(err, "Failed to login", map[string]interface{}{
			"email": loginRequest.Email,
		})
		return utils.UnauthorizedResponse(c, "Invalid email or password", nil)
	}

	logger.LogInfo("User logged in", map[string]interface{}{
		"id":    userDataOnLogin.UserID,
		"email": loginRequest.Email,
	})

	// Access token (short expiration)
	accessToken, err := utils.GenerateToken(fmt.Sprintf("%d", userDataOnLogin.UserID), userDataOnLogin.UserUUID, userDataOnLogin.Username, userDataOnLogin.RoleCode)
	if err != nil {
		logger.LogError(err, "Failed to generate access token", map[string]interface{}{
			"user_id": userDataOnLogin.UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	// Refresh token (long expiration)
	refreshToken, err := utils.GenerateRefreshToken(fmt.Sprintf("%d", userDataOnLogin.UserID), userDataOnLogin.UserUUID, userDataOnLogin.Username, userDataOnLogin.RoleCode)
	if err != nil {
		logger.LogError(err, "Failed to generate refresh token", map[string]interface{}{
			"user_id": userDataOnLogin.UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	// Save refresh token in the database
	err = utils.SaveRefreshToken(userDataOnLogin.UserUUID, refreshToken)
	if err != nil {
		logger.LogError(err, "Failed to save refresh token", map[string]interface{}{
			"user_id": userDataOnLogin.UserID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	responseData := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return utils.SuccessResponse(c, "User logged in successfully", responseData)
}

func (handler *authHandler) Logout(c *fiber.Ctx) error {
	log.Println("Logout function triggered")
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		log.Println("UserUUID not found in context")
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}
	log.Printf("UserUUID retrieved: %s\n", userUUID)

	// Delete WebSocket connection if exists
	conn, exists := utils.GetConnection(userUUID)
	if exists {
		log.Println("WebSocket connection exists, closing connection...")
		conn.Close()
		utils.RemoveConnection(userUUID)
		log.Printf("WebSocket connection for user %s closed and removed\n", userUUID)
	}

	err := handler.authService.DeleteRefreshTokenOnLogout(c.Context(), userUUID)
	if err != nil {
		log.Printf("Failed to delete refresh token for user %s: %v\n", userUUID, err)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	log.Printf("Refresh token for user %s deleted\n", userUUID)

	utils.InvalidateToken(c.Get("Authorization"))
	log.Println("Access token invalidated")

	err = handler.authService.UpdateUserStatus(userUUID, "offline", time.Now())
	if err != nil {
		log.Printf("Failed to update user status for user %s: %v\n", userUUID, err)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	log.Printf("User %s status updated to offline\n", userUUID)

	log.Println("Logout process completed successfully")
	return utils.SuccessResponse(c, "User logged out successfully", nil)
}

func (handler *authHandler) GetMyProfile(c *fiber.Ctx) error {
	log.Println("GetMyProfile function triggered")
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		log.Println("UserUUID not found in context")
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}
	log.Printf("UserUUID retrieved: %s\n", userUUID)

	roleCode, ok := c.Locals("role_code").(string)
	if !ok {
		log.Println("RoleCode not found in context")
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}
	log.Printf("RoleCode retrieved: %s\n", roleCode)

	user, err := handler.authService.GetMyProfile(userUUID, roleCode)
	if err != nil {
		log.Printf("Failed to retrieve profile for user %s: %v\n", userUUID, err)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	log.Printf("User profile retrieved for user %s\n", userUUID)

	return utils.SuccessResponse(c, "User profile retrieved", user)
}

// Reissue a new access token
func (handler *authHandler) IssueNewAccessToken(c *fiber.Ctx) error {
	refreshToken := c.Get("Authorization")
	if refreshToken == "" {
		return utils.UnauthorizedResponse(c, "Missing refresh token", nil)
	}

	// Remove "Bearer " prefix
	const bearerPrefix = "Bearer "
	if len(refreshToken) > len(bearerPrefix) && refreshToken[:len(bearerPrefix)] == bearerPrefix {
		refreshToken = refreshToken[len(bearerPrefix):]
	}

	claims, err := utils.ValidateToken(refreshToken)
	if err != nil {
		logger.LogWarn("Invalid refresh token", map[string]interface{}{
			"error": err.Error(),
		})
		return utils.UnauthorizedResponse(c, "Invalid refresh token", nil)
	}

	userID := claims["sub"].(string)
	userUUID := claims["user_uuid"].(string)

	tokenErr := handler.authService.CheckStoredRefreshToken(userUUID, refreshToken)
	if tokenErr != nil {
		logger.LogError(tokenErr, "Failed to get stored refresh token", map[string]interface{}{
			"user_id": userID,
			"token":   refreshToken,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	username := claims["user_name"].(string)
	roleCode := claims["role_code"].(string)

	err = handler.authService.UpdateRefreshToken(userUUID, refreshToken)
	if err != nil {
		logger.LogError(err, "Failed to update refresh token", map[string]interface{}{
			"user_uuid": userUUID,
		})
		return utils.UnauthorizedResponse(c, "Your session has expired or revoked, please login again", nil)
	}

	// Generate new access token
	accessToken, err := utils.GenerateToken(userID, userUUID, username, roleCode)
	if err != nil {
		logger.LogError(err, "Failed to generate access token", map[string]interface{}{
			"user_id": userID,
		})
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Access token refreshed", map[string]interface{}{
		"reissued_access_token": accessToken,
	})
}

func (handler *authHandler) AddDeviceToken(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		return utils.UnauthorizedResponse(c, "Token is invalid", nil)
	}

	tokenRequest := new(dto.DeviceTokenRequest)
	if err := c.BodyParser(tokenRequest); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	deviceToken := tokenRequest.Token

	// Save Device Token
	err := handler.authService.AddDeviceToken(userUUID, deviceToken)
	if err != nil {
		logger.LogError(err, "Failed to save FCM token", map[string]interface{}{
			"user_uuid":    userUUID,
			"device_token": tokenRequest.Token,
		})
		return utils.InternalServerErrorResponse(c, "Failed to save FCM token", nil)
	}

	return utils.SuccessResponse(c, "Device token added successfully", nil)
}
