package routes

import (
	"enterprise-helpdesk/internal/transport/http/handlers"
	"enterprise-helpdesk/internal/transport/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(router fiber.Router, deps Dependencies) {
	h := handlers.NewAuthHandler(deps)

	auth := router.Group("/auth")
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)

	auth.Use(middleware.AuthRequired(deps))

	auth.Get("/me", h.Me)
	auth.Get("/csrf", h.CSRF)
	auth.Get("/logout", h.Logout)
}
