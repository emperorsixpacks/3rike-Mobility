package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// Auth validates the Bearer JWT and checks the session is still alive in Redis.
func Auth(jwtSecret string, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid claims"})
		}

		sessionID, _ := claims["session_id"].(string)

		// Verify session is alive in Redis (skip check if Redis unavailable)
		if rdb != nil && sessionID != "" {
			if err := rdb.Get(c.Context(), "session:"+sessionID).Err(); err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session expired or logged out"})
			}
		}

		c.Locals("userID", uint(claims["sub"].(float64)))
		c.Locals("role", claims["role"].(string))
		c.Locals("sessionID", sessionID)
		return c.Next()
	}
}

// RequireRole restricts access to specific roles.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		if !allowed[role] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}
		return c.Next()
	}
}
