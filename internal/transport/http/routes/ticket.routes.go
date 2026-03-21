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
	tickets.Get("/:id", h.getbyid)

	tickets.Post("/", h.create)

	tickets.Patch("/:id", h.update)

	tickets.Get("/:id/comments", h.getcomments)

	tickets.Post("/:id/comments", h.addcomment)
}
