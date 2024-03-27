package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/sb-projects/sb-common/logger"
)

/*
func main() {
	ctx := context.Background()
	logger := logger.NewLogger()
	addr := os.Getenv("KAFKA_ADD")

	logger.Info(ctx, "Kafka topic creation -start")

	logger.Info(ctx, fmt.Sprintf("dialing %q", addr))
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		logger.Error(ctx, "failed to get connection", slog.Any("error", err))
		return
	}
	//defer conn.Close()

	logger.Info(ctx, "Connection established", slog.Group("connection",
		slog.String("host", conn.Broker().Host), slog.Int("id", conn.Broker().ID),
		slog.String("rack", conn.Broker().Rack), slog.Int("port", conn.Broker().Port)))

	logger.Info(ctx, "Identifying leader")
	controller, err := conn.Controller()
	if err != nil {
		logger.Error(ctx, "failed to get leader properties", slog.Any("error", err))
		return
	}

	logger.Info(ctx, fmt.Sprintf("leader address identified at %q", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))))

	var connLeader *kafka.Conn

	if addr == net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)) {
		connLeader = conn
	} else {
		connLeader, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
		if err != nil {
			logger.Error(ctx, "failed to connect leader", slog.Any("error", err))
			return
		}
	}
	defer connLeader.Close()
	logger.Info(ctx, "Leader Connection established", slog.Group("leader",
		slog.String("host", connLeader.Broker().Host), slog.Int("id", connLeader.Broker().ID),
		slog.String("rack", connLeader.Broker().Rack), slog.Int("port", connLeader.Broker().Port)))

	logger.Info(ctx, "Creating topics")

	partitions, err := connLeader.ReadPartitions()
	if err != nil {
		logger.Error(ctx, "failed to read partitions", slog.Any("error", err))
		return
	}

	printPartitionDetails(ctx, logger, partitions)

	topics := os.Getenv("TOPICS")
	if topics == "" {
		logger.Error(ctx, "No topics were set for TOPICS env variable", slog.String("topics", topics))
		return
	}

	topicConfigs := make([]kafka.TopicConfig, 0)

	for _, v := range strings.Split(topics, ",") {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             v,
			NumPartitions:     100,
			ReplicationFactor: 3,
		})
	}

	err = connLeader.CreateTopics(topicConfigs...)
	if err != nil {
		logger.Error(ctx, "failed to create topics", slog.Any("error", err))
		return
	}

	partitions, err = connLeader.ReadPartitions()
	if err != nil {
		logger.Error(ctx, "failed to read partitions", slog.Any("error", err))
		return
	}

	printPartitionDetails(ctx, logger, partitions)

	logger.Info(ctx, "Kafka topic creation -end")
}

func printPartitionDetails(ctx context.Context, logger *logger.Logger, partitions []kafka.Partition) {
	logger.Info(ctx, fmt.Sprintf("Total %d Partitions found ", len(partitions)))

	m := map[string]int{}

	for _, p := range partitions {
		m[p.Topic] = p.ID
	}
	for k, v := range m {
		logger.Info(ctx, "Topic", slog.String("name", k), slog.Int("partition", v))
	}
	logger.Info(ctx, fmt.Sprintf("Print completed for %d topics", len(m)))
}
*/

func main() {
	ctx := context.Background()
	logger := logger.NewLogger()

	brokers, topics, err := getInitData()
	if err != nil {
		logger.Error(ctx, "Topics were not created", slog.Any("error", err))
	}

	logFields := []any{slog.String("brokers", brokers)}
	for i, v := range topics {
		logFields = append(logFields, slog.String(fmt.Sprintf("topic-%d", i), v))
	}
	logger.Info(ctx, "Received request to create topics", slog.Group("kafka-config", logFields...))

	admin, err := createAdmin(brokers)
	if err != nil {
		log.Fatalf("Error creating admin: %v", err)
	}
	defer admin.Close()

	for _, topic := range topics {
		err := createTopic(admin, topic)
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("Error creating topic %s", topic), slog.Any("error", err))
		} else {
			logger.Debug(ctx, fmt.Sprintf("Created topic %s", topic))
		}
	}
}

func getInitData() (string, []string, error) {
	create := os.Getenv("CREATE_TOPICS")
	if create == "no" {
		return "", nil, fmt.Errorf("env CREATE_TOPICS is set to no")
	}

	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return "", nil, fmt.Errorf("env CREATE_TOPICS is set to no")
	}
	topics := os.Getenv("TOPICS")
	if brokers == "" {
		return "", nil, fmt.Errorf("env CREATE_TOPICS is set to no")
	}

	return brokers, strings.Split(topics, ","), nil
}

func createAdmin(brokerAddresses string) (sarama.ClusterAdmin, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0 // Use the desired Kafka version

	admin, err := sarama.NewClusterAdmin(strings.Split(brokerAddresses, ","), config)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func createTopic(admin sarama.ClusterAdmin, topic string) error {
	// Get the list of existing topics
	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}

	// Check if the topic already exists
	if _, exists := topics[topic]; exists {
		return fmt.Errorf("topic %s already exists", topic)
	}

	// If the topic doesn't exist, create it
	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1, // Set the number of partitions
		ReplicationFactor: 3, // Set the replication factor
	}, false)
	if err != nil {
		return err
	}

	return nil
}
