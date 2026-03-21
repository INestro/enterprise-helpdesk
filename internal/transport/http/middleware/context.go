package middleware

import "github.com/gofiber/fiber/v2"

type UserContext struct {
	UserID string
	Roles  []string
}

func GetUser(c *fiber.Ctx) *UserContext {
	user, _ := c.Locals("user").(*UserContext)
	return user
}
