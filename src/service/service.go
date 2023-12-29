package service

import (
	"github.com/sb-projects/sb-login/src/client"
)

type (
	Service interface {
	}
	impl struct {
		srvClient client.Client
	}
)
