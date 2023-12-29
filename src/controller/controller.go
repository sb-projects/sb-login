package controller

import (
	"github.com/go-playground/validator/v10"
)

type (
	Controller interface {
	}
	impl struct {
		validator *validator.Validate
	}
)

func NewController() Controller {
	return &impl{}
}
