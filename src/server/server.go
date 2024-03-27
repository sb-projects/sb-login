package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/sb-projects/sb-common/logger"
	ctxUtil "github.com/sb-projects/sb-common/util/ctx"
	"github.com/sb-projects/sb-login/src/controller"
	"github.com/sb-projects/sb-login/src/service"
	"github.com/sb-projects/sb-login/src/util/config"
)

func Serve() error {
	appConfig := config.Load()
	ctx := context.Background()
	ctx = ctxUtil.SetRequestID(ctx, "request-id")
	ctx = ctxUtil.SetTransactionID(ctx, "transaction-id")

	serviceLogger := logger.NewLogger()
	serviceLogger.Info(ctx, "Starting service", slog.String("Service", appConfig.App.Name))

	svr, err := service.NewService(ctx, serviceLogger, appConfig)
	if err != nil {
		serviceLogger.Error(ctx, "Failed to create Controller", slog.Any("err", err))
		return err
	}

	err = svr.Healthy()
	if err != nil {
		serviceLogger.Error(ctx, "Server not healthy", slog.Any("err", err))
		return err
	}

	ctrl, err := controller.NewController(serviceLogger, svr)
	if err != nil {
		serviceLogger.Error(ctx, "Failed to create Controller", slog.Any("err", err))
		return err
	}

	defer svr.Close(context.TODO())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.App.Port),
		Handler: ctrl.Routes(appConfig.App.Name),
	}

	// shutdownErrorChan := make(chan error)

	// go func() {
	// 	quitChan := make(chan os.Signal, 1)
	// 	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, os.Interrupt)

	// 	<-quitChan
	// 	serviceLogger.Info(ctx, "Service stopped 1")
	// 	shutdownErrorChan <- nil
	// }()

	// err := server.ListenAndServe()
	// if !errors.Is(err, http.ErrServerClosed) {
	// 	return err
	// }
	// serviceLogger.Info(ctx, "Started stopped")
	// err = <-shutdownErrorChan
	// if err != nil {
	// 	serviceLogger.Error(ctx, "Stopping service")
	// 	return err
	// }
	serviceLogger.Info(ctx, "Starting service", slog.Group("service",
		slog.String("name", appConfig.App.Name), slog.String("port", appConfig.App.Port)))
	log.Fatal(server.ListenAndServe())

	return nil
}
