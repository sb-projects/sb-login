package messaging

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

/*func connect(ctx context.Context, logger *logger.Logger) (*kafka.Conn, error) {

	addr := os.Getenv("KAFKA_BASE_ADD")
	if len(addr) == 0 {
		return nil, fmt.Errorf("empty kafka address")
	}
	topic := os.Getenv("KA_AUDIT")
	if len(addr) == 0 {
		return nil, fmt.Errorf("empty topic name")
	}
	logger.Debug(ctx, "Connecting to add", slog.String("addr", addr))

	var cancelF context.CancelFunc
	ctx, cancelF = context.WithTimeout(ctx, 10*time.Second)

	defer cancelF()

	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to dial %q", addr), err)
	}
	logger.Debug(ctx, "Dialled connection", slog.String("addr", addr))

	var broker kafka.Broker
	broker, err = conn.Controller()
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to get broker"), err)
	}

	addr = broker.Host + ":" + strconv.Itoa(broker.Port)
	conn, err = kafka.DialLeader(ctx, "tcp", addr, topic, 100)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to dial leader at %q", addr), err)
	}

	return conn, nil
}
*/

func connect() (sarama.SyncProducer, error) {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		return nil, fmt.Errorf("env KAFKA_BROKERS is empty")
	}

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
