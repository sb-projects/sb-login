package controller

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/sb-projects/sb-login/src/models"
	"github.com/sb-projects/sb-login/src/pkg/request"
	"github.com/sb-projects/sb-login/src/pkg/response"
)

func (c *Controller) registerV1(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		reqBody = models.RegisterUserReq{}
		err     error
	)
	c.logger.Debug(ctx, fmt.Sprintf("%q endpoint hit", r.URL))

	// Decode req
	err = request.DecodeJSON(w, r, &reqBody)
	if err != nil {
		c.logger.Error(ctx, "failed to decode request", slog.String("error", err.Error()))
		response.JSONError(w, http.StatusBadRequest, err)
		return
	}

	// register user
	err = c.srv.RegisterUser(ctx, reqBody)
	if err != nil {
		c.logger.Error(ctx, "failed to register user", slog.String("error", err.Error()))
		response.JSONError(w, http.StatusInternalServerError, err)
		return
	}

	data := map[string]string{
		"Status": "OK",
	}

	response.JSON(w, http.StatusCreated, data)
}
