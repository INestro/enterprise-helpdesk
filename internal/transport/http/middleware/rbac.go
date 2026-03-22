package middleware

import "github.com/gofiber/fiber/v2"

func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCtx, ok := c.Locals("user").(*UserContext)
		if !ok {
			return fiber.ErrUnauthorized
		}

		for _, r := range userCtx.Roles {
			if r == role {
				return c.Next()
			}
		}

		return fiber.ErrForbidden
	}
}
