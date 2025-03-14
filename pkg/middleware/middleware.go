package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rk-the-dev/golib-core/pkg/logger"
	"github.com/sirupsen/logrus"
)

// JWTMiddleware verifies the JWT token
func JWTMiddleware(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}
		// Token is valid, store claims in context
		c.Locals("user", token.Claims)
		return c.Next()
	}
}

// LoggingMiddleware logs incoming HTTP requests
func LoggingMiddleware(c *fiber.Ctx) error {
	logger.Info("HTTP Request", logrus.Fields{
		"method": c.Method(),
		"path":   c.Path(),
	})
	return c.Next()
}

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware(c *fiber.Ctx) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Recovered from panic", logrus.Fields{
				"error": err,
				"path":  c.Path(),
			})
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}
	}()
	return c.Next()
}

// CORSMiddleware sets CORS headers
func CORSMiddleware(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	if c.Method() == fiber.MethodOptions {
		return c.SendStatus(fiber.StatusNoContent)
	}
	return c.Next()
}
