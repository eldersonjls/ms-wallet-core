package create_client

import (
	"context"
	"time"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
	"github.com/eldersonjls/ms-wallet-core/internal/gateway"
	"github.com/eldersonjls/ms-wallet-core/pkg/events"
	"github.com/eldersonjls/ms-wallet-core/pkg/uow"
)

type CreateClientInputDTO struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	BankID      string    `json:"bank_id"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

type CreateClientOutputDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	BankID      string    `json:"bank_id"`
	DateOfBirth time.Time `json:"date_of_birth"`
	CreatedAt   time.Time `json:"create_at"`
}

type CreateClientUseCase struct {
	Uow             uow.UowInterface
	EventDispatcher events.EventDispatcherInterface
	ClientCreated   events.EventInterface
}

func NewCreateClientUseCase(
	Uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	clientCreated events.EventInterface,
) *CreateClientUseCase {
	return &CreateClientUseCase{
		Uow:             Uow,
		EventDispatcher: eventDispatcher,
		ClientCreated:   clientCreated,
	}
}

func (uc *CreateClientUseCase) Execute(ctx context.Context, input CreateClientInputDTO) (*CreateClientOutputDTO, error) {
	output := &CreateClientOutputDTO{}
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		clientRepository := uc.getClientRepository(ctx)
		client, err := entity.NewClient(input.Name, input.Email, input.BankID, input.DateOfBirth)
		if err != nil {
			return err
		}
		err = clientRepository.Save(client)
		if err != nil {
			return err
		}
		output.ID = client.ID
		output.Name = client.Name
		output.Email = client.Email
		output.BankID = client.BankID
		output.DateOfBirth = client.DateOfBirth
		output.CreatedAt = client.CreatedAt
		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.ClientCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.ClientCreated)

	return output, nil
}

func (uc *CreateClientUseCase) getClientRepository(ctx context.Context) gateway.ClientGateway {
	repo, err := uc.Uow.GetRepository(ctx, "ClientDB")
	if err != nil {
		panic(err)
	}

	if repo == nil {
		panic("ClientRepository is nil")
	}

	return repo.(gateway.ClientGateway)
}
