package create_transaction

import (
	"context"
	"testing"
	"time"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
	"github.com/eldersonjls/ms-wallet-core/internal/event"
	"github.com/eldersonjls/ms-wallet-core/internal/usecase/mocks"
	"github.com/eldersonjls/ms-wallet-core/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client1, _ := entity.NewClient("John Doe", "john@example.com", "123", clientDateOfBirth)
	account1, _ := entity.NewAccount(client1, entity.CurrentAccount)

	account1.Credit(1000)

	client2, _ := entity.NewClient("Jane Doe", "to@example.com", "1234", clientDateOfBirth)
	account2, _ := entity.NewAccount(client2, entity.CurrentAccount)
	account2.Credit(1000)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: account1.ID,
		AccountIDTo:   account2.ID,
		Amount:        100,
	}

	dispatcher := events.NewEventDispatcher()
	eventTransaction := event.NewTransactionCreated()
	eventBalance := event.NewBalanceUpdated()

	ctx := context.Background()

	uc := NewCreateTransactionUseCase(mockUow, dispatcher, eventTransaction, eventBalance)
	output, err := uc.Execute(ctx, inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)

}
