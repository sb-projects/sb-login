package controller

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (c *Controller) Routes(appName string) http.Handler {
	mux := mux.NewRouter()

	baseRoute := mux.PathPrefix("/" + appName).Subrouter()

	intRoute := baseRoute.PathPrefix("/internal").Subrouter()
	intRoute.HandleFunc("/health", c.health).Methods(http.MethodGet)

	v1routes := baseRoute.PathPrefix("/v1").Subrouter()
	v1routes.HandleFunc("/register", c.registerV1).Methods(http.MethodPost)
	v1routes.HandleFunc("/login", c.loginV1).Methods(http.MethodPost)

	return mux
}
