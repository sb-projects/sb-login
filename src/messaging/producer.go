package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/sb-projects/sb-common/logger"
	"github.com/sb-projects/sb-login/src/models"
)

type (
	Producer interface {
		// Send(context.Context, any) error
		Healthy() error
		Close() error
		Write(context.Context, models.KafkaWriteObj) error
	}
	impl struct {
		logger *logger.Logger
		prod   sarama.SyncProducer
	}
)

// Healthy implements Producer.
func (i *impl) Healthy() error {
	isTransactional := i.prod.IsTransactional()
	if !isTransactional {
		return errors.New("producer not transactional")
	}
	return nil
}

// Close implements Producer.
func (i *impl) Close() error {
	err := i.prod.Close()
	return err
}

func New(ctx context.Context, logger *logger.Logger) (Producer, error) {
	leader, err := connect()
	if err != nil {
		logger.Debug(ctx, "failed to get producer ", slog.Any("error", err))
		return nil, err
	}
	logger.Debug(ctx, "leader connected")
	return &impl{
		logger: logger,
		prod:   leader,
	}, nil
}

func (i *impl) Write(ctx context.Context, messageObj models.KafkaWriteObj) error {

	messageVal, err := json.Marshal(messageObj.Message)
	if err != nil {
		i.logger.Error(ctx, "Failed to marshal message data", slog.Any("error", err))
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     messageObj.Topic,
		Value:     sarama.StringEncoder(messageVal),
		Timestamp: time.Now().UTC(),
	}

	partition, offset, err := i.prod.SendMessage(msg)
	if err != nil {
		i.logger.Error(ctx, "Failed to Send message", slog.Any("error", err))
		return err
	}

	i.logger.Debug(ctx, "message sent", slog.Group("message", slog.String("topic", messageObj.Topic),
		slog.Int("partition", int(partition)), slog.Int64("offset", offset)))
	return nil

}
