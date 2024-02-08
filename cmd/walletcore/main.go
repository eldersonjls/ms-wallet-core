package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/eldersonjls/ms-wallet-core/internal/database"
	"github.com/eldersonjls/ms-wallet-core/internal/event"
	"github.com/eldersonjls/ms-wallet-core/internal/event/handler"

	"github.com/eldersonjls/ms-wallet-core/internal/usecase/create_account"
	"github.com/eldersonjls/ms-wallet-core/internal/usecase/create_client"
	"github.com/eldersonjls/ms-wallet-core/internal/usecase/create_transaction"

	"github.com/eldersonjls/ms-wallet-core/internal/web"
	"github.com/eldersonjls/ms-wallet-core/internal/web/webserver"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/eldersonjls/ms-wallet-core/pkg/events"
	"github.com/eldersonjls/ms-wallet-core/pkg/kafka"
	"github.com/eldersonjls/ms-wallet-core/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()

	eventDispatcher.Register("ClientCreated", handler.NewClientCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("AccountCreated", handler.NewAccountCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))

	clientCreatedEvent := event.NewClientCreated()
	accountCreatedEvent := event.NewAccountCreated()
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("ClientDB", func(tx *sql.Tx) interface{} {
		return database.NewClientDB(db)
	})

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createClientUseCase := create_client.NewCreateClientUseCase(uow, eventDispatcher, clientCreatedEvent)
	createAccountUseCase := create_account.NewCreateAccountUseCase(uow, eventDispatcher, accountCreatedEvent)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")
	webserver.Start()

}
