package handler

import (
	"fmt"
	"sync"

	"github.com/eldersonjls/ms-wallet-core/pkg/events"
	"github.com/eldersonjls/ms-wallet-core/pkg/kafka"
)

type ClientCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewClientCreatedKafkaHandler(kafka *kafka.Producer) *ClientCreatedKafkaHandler {
	return &ClientCreatedKafkaHandler{
		Kafka: kafka,
	}
}

func (h *ClientCreatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(message, nil, "clients")
	fmt.Println("ClientCreatedKafkaHandler: ", message.GetPayload())
}
