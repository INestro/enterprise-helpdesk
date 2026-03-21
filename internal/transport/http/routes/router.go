package routes

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App, deps Dependencies) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	registerHealthRoutes(v1)
}

func registerHealthRoutes(router fiber.Router) {
	router.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
