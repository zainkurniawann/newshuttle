package middleware

import (
	"shuttle/logger"
	"shuttle/utils"
	"shuttle/services"

	"github.com/gofiber/fiber/v2"
)

func SchoolAdminMiddleware(service services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userUUID, ok := c.Locals("userUUID").(string)
		if !ok || userUUID == "" {
			return utils.UnauthorizedResponse(c, "User ID is missing or invalid", nil)
		}

		schoolUUID, err := service.CheckPermittedSchoolAccess(userUUID)
		if err != nil {
			return utils.ForbiddenResponse(c, "You don't have permission to any school, please contact the support team", nil)
		}

		c.Locals("schoolUUID", schoolUUID)

		return c.Next()
	}
}

func AuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return utils.UnauthorizedResponse(c, "Missing token", nil)
		}

		const bearerPrefix = "Bearer "
		if len(token) > len(bearerPrefix) && token[:len(bearerPrefix)] == bearerPrefix {
			token = token[len(bearerPrefix):]
		}

		_, exists := utils.InvalidTokens[token]
		if exists {
			return utils.UnauthorizedResponse(c, "Invalid token or you have been logged out", nil)
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			logger.LogWarn("Invalid token", map[string]interface{}{"error": err.Error()})
			return utils.UnauthorizedResponse(c, "Token is invalid", nil)
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			logger.LogWarn("User ID is missing or invalid", map[string]interface{}{"claims": claims})
			return utils.UnauthorizedResponse(c, "Token is invalid", nil)
		}

		userUUID, ok := claims["user_uuid"].(string)
		if !ok || userUUID == "" {
			logger.LogWarn("User UUID is missing or invalid", map[string]interface{}{"claims": claims})
			return utils.UnauthorizedResponse(c, "Token is invalid", nil)
		}

		role_code, ok := claims["role_code"].(string)
		if !ok || role_code == "" {
			logger.LogWarn("Role code is missing or invalid", map[string]interface{}{"claims": claims})
			return utils.UnauthorizedResponse(c, "Token is invalid", nil)
		}

		user_name, ok := claims["user_name"].(string)
		if !ok || user_name == "" {
			logger.LogWarn("User name is missing or invalid", map[string]interface{}{"claims": claims})
			return utils.UnauthorizedResponse(c, "Token is invalid", nil)
		}

		c.Locals("userID", userID)
		c.Locals("userUUID", userUUID)
		c.Locals("role_code", role_code)
		c.Locals("user_name", user_name)

		return c.Next()
	}
}

func AuthorizationMiddleware(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role_code, ok := c.Locals("role_code").(string)
		if !ok || role_code == "" {
			return utils.UnauthorizedResponse(c, "Role code is missing or invalid", nil)
		}

		if !contains(allowedRoles, role_code) {
			return utils.ForbiddenResponse(c, "You don't have permission to access this resource", nil)
		}

		return c.Next()
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
