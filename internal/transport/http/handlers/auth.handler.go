package handlers

import (
	"enterprise-helpdesk/internal/transport/http/handlers/dto"
	"enterprise-helpdesk/internal/transport/http/routes"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	deps routes.Dependencies
}

func NewAuthHandler(deps routes.Dependencies) *AuthHandler {
	return &AuthHandler{deps: deps}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.Email != "admin@example.com" || req.Password != "password" {
		return fiber.ErrUnauthorized
	}

	user := fiber.Map{
		"id":        "user-1",
		"email":     req.Email,
		"full_name": "Admin User",
		"roles":     []string{"admin"},
	}

	token, err := h.gener
}
