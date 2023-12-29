package models

type (
	AppConfig struct {
		Name string `validate:"required"`
		Port string `validate:"required"`
	}
	DBConfig struct {
		Name string `validate:"required"`
		URL  string `validate:"required"`
	}

	Config struct {
		App AppConfig
		DB  DBConfig
	}
)
