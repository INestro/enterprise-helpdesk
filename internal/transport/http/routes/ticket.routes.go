package routes

import (
	"enterprise-helpdesk/internal/transport/http/handlers"
	"enterprise-helpdesk/internal/transport/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func registerTicketRoutes(router fiber.Router, deps Dependencies) {
	h := handlers.NewTicketHandler(deps)

	tickets := router.Group("/tickets",
		middleware.AuthRequired(deps),
	)

	tickets.Get("/", h.GetList)
	tickets.Get("/:id", h.GetByID)

	tickets.Post("/",
		middleware.CSRF(deps),
		middleware.Idempotency(deps),
		h.Create,
	)

	tickets.Patch("/:id",
		middleware.CSRF(deps),
		h.Update,
	)

	tickets.Get("/:id/comments", h.GetComments)

	tickets.Post("/:id/comments",
		middleware.CSRF(deps),
		h.AddComement)
}
