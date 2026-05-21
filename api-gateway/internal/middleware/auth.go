package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "aitu-superapp-secret-2026"

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Поддержка старого demo-токена
		if tokenStr == "super-secret-token" {
			c.Request().Header.Set("X-User-ID", "123e4567-e89b-12d3-a456-426614174000")
			c.Request().Header.Set("X-User-Role", "student")
			return c.Next()
		}

		// Валидация настоящего JWT
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})
		}

		userID, _ := claims["user_id"].(string)
		userRole, _ := claims["role"].(string)
		c.Request().Header.Set("X-User-ID", userID)
		c.Request().Header.Set("X-User-Role", userRole)
		return c.Next()
	}
}
