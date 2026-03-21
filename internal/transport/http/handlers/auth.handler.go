package handlers

import (
	"enterprise-helpdesk/internal/transport/http/dto"
	"enterprise-helpdesk/internal/transport/http/middleware"
	"enterprise-helpdesk/internal/transport/http/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

	token, err := h.generateJWT(user["id"].(string), user["roles"].([]string))
	if err != nil {
		return err
	}

	// refresh token кладём в cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,
	})

	return c.JSON(dto.AuthResponse{
		AccessToken: token,
		User:        user,
	})
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.ErrUnauthorized
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.deps.Cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}

	claims := token.Claims.(jwt.MapClaims)

	userID, _ := claims["user_id"].(string)
	rolesRaw, _ := claims["roles"].([]interface{})

	var roles []string
	for _, r := range rolesRaw {
		if s, ok := r.(string); ok {
			roles = append(roles, s)
		}
	}

	newToken, err := h.generateJWT(userID, roles)
	if err != nil {
		return err
	}

	return c.JSON(dto.RefreshResponse{
		AccessToken: newToken,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := middleware.GetUser(c)

	return c.JSON(fiber.Map{
		"id":    user.UserID,
		"roles": user.Roles,
	})
}

func (h *AuthHandler) CSRF(c *fiber.Ctx) error {
	token := generateCSRFToken()

	c.Cookie(&fiber.Cookie{
		Name:  "csrf_token",
		Value: token,
	})

	return c.JSON(dto.CSRFResponse{
		CSRFToken: token,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("refresh_token")
	c.ClearCookie("csrf_token")

	return c.JSON(fiber.Map{
		"message": "logout",
	})
}

func (h *AuthHandler) generateJWT(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"roles":   roles,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(h.deps.Cfg.JWTSecret))
}

func generateCSRFToken() string {
	return time.Now().Format("20060102150405")
}
