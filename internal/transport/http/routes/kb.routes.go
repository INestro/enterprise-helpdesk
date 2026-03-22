package routes

import (
	"enterprise-helpdesk/internal/transport/http/handlers"
	"enterprise-helpdesk/internal/transport/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func registerKnowledgeBaseRoutes(router fiber.Router, deps Dependencies) {
	h := handlers.NewKnowledgeBaseHandler(deps)

	kb := router.Group("knowledge",
		middleware.AuthRequired(deps),
	)

	kb.Get("/articles", h.getlist)
	kb.Get("/articles/:id", h.getbyid)

	kb.Post("/articles",
		middleware.CSRF(deps),
		middleware.RequireRole("admin"),
		h.create,
	)
}
