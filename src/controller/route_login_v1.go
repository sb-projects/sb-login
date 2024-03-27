package controller

import (
	"fmt"
	"net/http"

	"github.com/sb-projects/sb-login/src/pkg/response"
)

func (c *Controller) loginV1(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)
	data := map[string]string{
		"Status": "OK",
	}
	c.logger.Info(ctx, fmt.Sprintf("%q endpoint hit", r.URL))
	response.JSON(w, http.StatusOK, data)
}
