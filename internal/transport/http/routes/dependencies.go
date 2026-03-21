package routes

import (
	"enterprise-helpdesk/internal/config"
	"enterprise-helpdesk/internal/infrastructure/database"
	"enterprise-helpdesk/internal/infrastructure/redis"
)

type Dependencies struct {
	DB    *database.Postgres
	Redis *redis.Client
	Cfg   *config.Config
}
