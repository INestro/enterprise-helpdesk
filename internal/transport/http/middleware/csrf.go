package middleware

import (
	"enterprise-helpdesk/internal/transport/http/routes"

	"github.com/gofiber/fiber/v2"
)

func CSRF(deps routes.Dependencies) fiber.Handler {
	return func(c *fiber.Ctx) error {
		method := c.Method()

		if method == fiber.MethodGet || method == fiber.MethodHead {
			return c.Next()
		}

		csrfHeader := c.Get("X-CSRF-Token")
		csrfCookie := c.Cookies("csrf_token")

		if csrfHeader == "" || csrfCookie == "" {
			return fiber.ErrForbidden
		}

		if csrfHeader != csrfCookie {
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
