package create_account

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

func TestCreateAccountUseCase_Execute(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := entity.NewClient("John Doe", "john@example.com", "123", clientDateOfBirth)

	mockUow := &mocks.UowMock{}
	mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	inputDto := CreateAccountInputDTO{
		ClientID:    client.ID,
		AccountType: entity.CurrentAccount,
	}

	dispatcher := events.NewEventDispatcher()
	eventAccountCreated := event.NewAccountCreated()

	ctx := context.Background()

	uc := NewCreateAccountUseCase(mockUow, dispatcher, eventAccountCreated)
	output, err := uc.Execute(ctx, inputDto)

	assert.Nil(t, err)
	assert.NotNil(t, output)

	mockUow.AssertExpectations(t)
	mockUow.AssertNumberOfCalls(t, "Do", 1)
}
