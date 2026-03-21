package app

import (
	"enterprise-helpdesk/internal/config"
	"enterprise-helpdesk/internal/infrastructure/database"
	"enterprise-helpdesk/internal/infrastructure/redis"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type App struct {
	fiber *fiber.App
	cfg   *config.Config
	db    *database.Postgres
	redis *redis.Client
}

func New() (*App, error) {
	cfg := config.Load()

	app := &App{
		cfg: cfg,
	}

	if err := app.initInfrastructure(); err != nil {
		return nil, err
	}

	app.initFiber()
	//app.
}

func (a *App) initInfrastructure() error {
	db, err := database.NewPostgres(a.cfg.PostgreDSN)
	if err != nil {
		return err
	}
	a.db = db

	rdb, err := redis.New(a.cfg.RedisAddr)
	if err != nil {
		return err
	}
	a.redis = rdb

	return nil
}

func (a *App) initFiber() {
	a.fiber = fiber.New(fiber.Config{
		AppName:      "Enterprise HelpDesk",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// panic recovery
	a.fiber.Use(recover.New())
}

func (a *App) initRoutes() {
	//rout
}
