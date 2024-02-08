package handler

import (
	"fmt"
	"sync"

	"github.com/eldersonjls/ms-wallet-core/pkg/events"
	"github.com/eldersonjls/ms-wallet-core/pkg/kafka"
)

type AccountCreatedKafkaHandler struct {
	Kafka *kafka.Producer
}

func NewAccountCreatedKafkaHandler(kafka *kafka.Producer) *AccountCreatedKafkaHandler {
	return &AccountCreatedKafkaHandler{
		Kafka: kafka,
	}
}

func (h *AccountCreatedKafkaHandler) Handle(message events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	h.Kafka.Publish(message, nil, "accounts")
	fmt.Println("AccountCreatedKafkaHandler: ", message.GetPayload())
}
