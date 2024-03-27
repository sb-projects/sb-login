package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/sb-projects/sb-common/logger"
	"github.com/sb-projects/sb-login/src/service"
)

type (
	Controller struct {
		logger    *logger.Logger
		validator *validator.Validate
		srv       service.Service
	}
)

func NewController(log *logger.Logger, service service.Service) (*Controller, error) {
	return &Controller{
		logger:    log,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		srv:       service,
	}, nil
}
