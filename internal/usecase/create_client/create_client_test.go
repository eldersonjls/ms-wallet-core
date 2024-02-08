package create_client

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

func TestCreateClientUseCase_Execute(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := entity.NewClient("client1", "j@j.com", "123", clientDateOfBirth)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	inputDto := CreateClientInputDTO{
		Name:        client.Name,
		Email:       client.Email,
		BankID:      client.BankID,
		DateOfBirth: client.DateOfBirth,
	}

	dispatcher := events.NewEventDispatcher()
	eventClientCreated := event.NewClientCreated()

	ctx := context.Background()

	uc := NewCreateClientUseCase(mockUow, dispatcher, eventClientCreated)
	output, err := uc.Execute(ctx, inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
