package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/sb-projects/sb-common/logger"
)

const (
	brokerAddressesEnvKey = "KAFKA_BROKERS" // Environment variable key for Kafka broker addresses
	topicNamesEnvKey      = "TOPICS"        // Environment variable key for topic names
)

func main() {
	ctx := context.Background()
	logger := logger.NewLogger()
	brokerAddresses := os.Getenv(brokerAddressesEnvKey)
	if brokerAddresses == "" {
		logger.Error(ctx, "KAFKA_BROKERS environment variable is not set")
		return
	}
	topicNamesStr := os.Getenv(topicNamesEnvKey)
	if topicNamesStr == "" {
		logger.Error(ctx, "TOPICS environment variable is not set")
		return
	}

	consumer, err := createConsumer(brokerAddresses)
	if err != nil {
		logger.Error(ctx, "Error creating consumer", slog.Any("error", err))
		return
	}
	defer consumer.Close()
	logger.Debug(ctx, "consumers started")

	topicNames := strings.Split(topicNamesStr, ",")
	var wg sync.WaitGroup

	for _, topic := range topicNames {
		wg.Add(1)
		go consumeMessages(consumer, topic, &wg, logger)
	}

	// Wait for a termination signal to gracefully stop the consumer
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	logger.Warn(ctx, "Received termination signal. Shutting down...")
	wg.Wait()
	logger.Info(ctx, "Consumer service stopped.")
}

func createConsumer(brokerAddresses string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(strings.Split(brokerAddresses, ","), config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func consumeMessages(consumer sarama.Consumer, topic string, wg *sync.WaitGroup, logger *logger.Logger) {
	defer wg.Done()
	ctx := context.Background()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		logger.Error(ctx, "Error getting partitions", slog.String("topic", topic), slog.Any("error", err))
		return
	}

	var partitionConsumers []sarama.PartitionConsumer

	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			logger.Warn(ctx, "Error creating partition consumer", slog.String("topic", topic), slog.Int("partition", int(partition)), slog.Any("error", err))
			continue
		}

		partitionConsumers = append(partitionConsumers, partitionConsumer)
	}

	for _, partitionConsumer := range partitionConsumers {
		go func(pc sarama.PartitionConsumer) {
			defer pc.AsyncClose()

			for {
				select {
				case msg := <-pc.Messages():
					logger.Info(ctx, "Received message", slog.Group("message-details", slog.String("topic", msg.Topic), slog.String("value", string(msg.Value)), slog.Int("partition", int(msg.Partition))))
					msgInd, _ := json.MarshalIndent(msg.Value, "", "\t")
					fmt.Println("Message start-\n", string(msgInd), "\n- Message end")
				case err := <-pc.Errors():
					logger.Error(ctx, "Error consuming message", slog.Any("error", err))
				}
			}
		}(partitionConsumer)
	}

	logger.Info(ctx, "Started consuming messages", slog.String("topic", topic))
	select {}
}
