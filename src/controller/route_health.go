package controller

import (
	"fmt"
	"net/http"

	"github.com/sb-projects/sb-login/src/pkg/response"
)

func (c *Controller) health(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	c.logger.Info(ctx, fmt.Sprintf("%q endpoint hit", r.URL))
	data := map[string]string{
		"Status": "OK",
	}

	response.JSON(w, http.StatusOK, data)
}
