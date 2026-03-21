package routes

import (
	"enterprise-helpdesk/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(router fiber.Router, deps Dependencies) {
	h := handlers.NewAuthHandler(deps)

	auth := router.Group("/auth")
	auth.Post("/login", h.)
}
