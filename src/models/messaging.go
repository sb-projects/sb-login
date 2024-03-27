package models

type (
	KafkaWriteObj struct {
		Topic     string
		Partition int
		Message   KafkaMessage
	}
	KafkaOperation struct {
		Object string `json:"object"`
		Type   string `json:"type"`
	}
	KafkaMessage struct {
		Version   string         `json:"version"`
		Operation KafkaOperation `json:"operation"`
		Data      map[string]any `json:"data"`
	}
)
