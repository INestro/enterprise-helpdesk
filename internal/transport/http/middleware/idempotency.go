package middleware

import (
	"context"
	"enterprise-helpdesk/internal/transport/http/routes"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Idempotency(deps routes.Dependencies) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("Idempotency-Key")
		if key == "" {
			return c.Next()
		}

		ctx := context.Background()
		redis := deps.Redis.RDB

		cacheKey := "idem:" + key

		exists, err := redis.Get(ctx, cacheKey).Result()
		if err == nil && exists != "" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "duplicate request",
			})
		}

		// помечаем как выполненный
		if err := redis.Set(ctx, cacheKey, "1", 10*time.Minute).Err(); err != nil {
			return err
		}

		return c.Next()
	}
}
