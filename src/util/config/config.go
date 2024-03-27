package config

import (
	"os"

	"github.com/sb-projects/sb-login/src/models"
)

const (
	appName string = "APP_NAME"
	appPort string = "APP_PORT"
	dbPGURL string = "APP_DB_PG_URL"
)

func Load() models.Config {
	config := models.Config{
		App: models.AppConfig{
			Name: os.Getenv(appName),
			Port: os.Getenv(appPort),
		},
		DB: models.DBConfig{
			Name:    "PG",
			URL:     os.Getenv(dbPGURL),
			Migrate: true,
		},
	}
	return config
}
