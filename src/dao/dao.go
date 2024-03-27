package dao

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sb-projects/sb-login/src/models"
	"github.com/sb-projects/sb-login/src/pkg/postgre"
)

type (
	Daolayer interface {
		RegisterUser(context.Context, models.RegisterUserReq) (string, error)
	}
	Dao struct {
		pg *sqlx.DB
	}
)

func New(config models.Config) (Daolayer, error) {
	dbPG, err := postgre.New(config.DB.URL, config.DB.Migrate)
	if err != nil {
		return nil, err
	}
	return &Dao{pg: dbPG}, nil
}
