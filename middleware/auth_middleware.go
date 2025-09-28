package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rdsarjito/marketplace-backend/constants"
	"github.com/rdsarjito/marketplace-backend/domain/dto/response"
	"github.com/rdsarjito/marketplace-backend/services"
	"github.com/rdsarjito/marketplace-backend/utils"
)

func AuthMiddleware(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(constants.ErrUnauthorized, nil))
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(constants.ErrUnauthorized, nil))
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(constants.ErrUnauthorized, nil))
		}

		// Validate token
		claims, err := utils.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(constants.ErrInvalidToken, nil))
		}

		// Check if user still exists
		user, err := userService.GetMyProfile(claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse(constants.ErrUserNotFound, nil))
		}

		// Set user data in context
		c.Locals("userID", claims.UserID)
		c.Locals("isAdmin", claims.IsAdmin)
		c.Locals("user", user)

		return c.Next()
	}
}

func AdminMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin := c.Locals("isAdmin").(bool)
		if !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(response.ErrorResponse(constants.ErrForbidden, nil))
		}
		return c.Next()
	}
}
