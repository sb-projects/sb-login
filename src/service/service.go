package service

import (
	"context"
	"log/slog"
	"os"

	"github.com/sb-projects/sb-common/logger"
	"github.com/sb-projects/sb-login/src/client"
	"github.com/sb-projects/sb-login/src/dao"
	"github.com/sb-projects/sb-login/src/messaging"
	"github.com/sb-projects/sb-login/src/models"
	"github.com/sb-projects/sb-login/src/pkg/secret"
)

type (
	Service interface {
		Healthy() error
		Close(context.Context) error
		RegisterUser(context.Context, models.RegisterUserReq) error
	}
	impl struct {
		log       *logger.Logger
		srvClient client.Client
		dao       dao.Daolayer
		producer  messaging.Producer
	}
)

// Healthy implements Service.
func (i *impl) Healthy() error {
	// err := i.producer.Healthy()
	// if err != nil {
	// 	return errors.Join(fmt.Errorf("producer not healthy", err))
	// }
	return nil
}

func NewService(ctx context.Context, log *logger.Logger, config models.Config) (Service, error) {
	dao, err := dao.New(config)
	if err != nil {
		log.Debug(ctx, "failed to connect with dao", slog.Any("error", err))
		return nil, err
	}
	log.Debug(ctx, "dao started")
	var prod messaging.Producer
	prod, err = messaging.New(ctx, log)
	if err != nil {
		log.Debug(ctx, "failed to connect with Kafka", slog.Any("error", err))
		return nil, err
	}
	log.Debug(ctx, "producer started started")
	service := impl{
		log:       log,
		srvClient: client.NewClient(),
		dao:       dao,
		producer:  prod,
	}
	return &service, nil
}

func (s *impl) Close(_ context.Context) error {
	s.log.Info(context.TODO(), "Closing server")
	return nil
}

// RegisterUser implements Service.
func (s *impl) RegisterUser(ctx context.Context, user models.RegisterUserReq) error {
	s.log.Debug(ctx, "RegisterUser: Started registration", slog.String("User", user.Email))

	var err error
	user, err = processUserReq(user)
	if err != nil {
		s.log.Error(ctx, "RegisterUser: Failed to pre-process request", slog.String("error", err.Error()))
		return err
	}

	var id string
	id, err = s.dao.RegisterUser(ctx, user)
	if err != nil {
		s.log.Error(ctx, "RegisterUser: Failed to register user", slog.String("error", err.Error()))
		return err
	}
	s.log.Info(ctx, "RegisterUser: Registered user", slog.String("User", id))
	err = s.producer.Write(ctx, models.KafkaWriteObj{
		Topic:     os.Getenv("KA_AUDIT"),
		Partition: 100,
		Message: models.KafkaMessage{
			Version: "1.0",
			Operation: models.KafkaOperation{
				Object: "user",
				Type:   "register",
			},
			Data: map[string]any{
				"userID": id,
				"name":   user.Name,
				"email":  user.Email,
			},
		},
	})
	if err != nil {
		s.log.Error(ctx, "RegisterUser: Failed to write message", slog.String("error", err.Error()))
	}
	return nil
}

func processUserReq(user models.RegisterUserReq) (models.RegisterUserReq, error) {
	hash, err := secret.Hash(user.HashPass)
	if err != nil {
		return user, err
	}
	user.HashPass = hash
	return user, nil
}
