package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/sb-projects/sb-common/logger"
)

func LogTxn(next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			method = r.Method
			url    = r.URL.String()
			proto  = r.Proto
		)

		requestAttrs := slog.Group("request", "method", method, "url", url, "proto", proto)

		logger.Info(context.TODO(), "access", requestAttrs)
		next.ServeHTTP(w, r)
	})
}
