package service

import (
	"github.com/sb-projects/sb-login/client"
)

type (
	Service interface {
	}
	impl struct {
		client client.Client
	}
)
