package middleware

import (
	"enterprise-helpdesk/internal/transport/http/routes"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(deps routes.Dependencies) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			return fiber.ErrUnauthorized
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(deps.Cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}

		userID, _ := claims["user_id"].(string)
		rolesRaw, _ := claims["roles"].([]interface{})

		var roles []string
		for _, r := range rolesRaw {
			if s, ok := r.(string); ok {
				roles = append(roles, s)
			}
		}

		ctx := &UserContext{
			UserID: userID,
			Roles:  roles,
		}

		c.Locals("user", ctx)

		return c.Next()
	}
}
