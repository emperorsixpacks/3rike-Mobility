package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// Cache returns a cache-aside middleware for GET requests.
// Skips caching if redis is nil (graceful degradation).
func Cache(rdb *redis.Client, ttl time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if rdb == nil || c.Method() != fiber.MethodGet {
			return c.Next()
		}
		key := "cache:" + c.Path() + "?" + string(c.Request().URI().QueryString())
		cached, err := rdb.Get(c.Context(), key).Bytes()
		if err == nil {
			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.Send(cached)
		}
		if err := c.Next(); err != nil {
			return err
		}
		if c.Response().StatusCode() == fiber.StatusOK {
			_ = rdb.Set(c.Context(), key, c.Response().Body(), ttl).Err()
		}
		return nil
	}
}
