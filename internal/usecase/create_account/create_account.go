package create_account

import (
	"context"
	"time"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
	"github.com/eldersonjls/ms-wallet-core/internal/gateway"
	"github.com/eldersonjls/ms-wallet-core/pkg/events"
	"github.com/eldersonjls/ms-wallet-core/pkg/uow"
)

type CreateAccountInputDTO struct {
	ClientID    string             `json:"client_id"`
	AccountType entity.AccountType `json:"account_type"`
}

type CreateAccountOutputDTO struct {
	ID          string             `json:"id"`
	ClientID    string             `json:"client_id"`
	AccountType entity.AccountType `json:"account_type"`
	CreatedAt   time.Time          `json:"create_at"`
}

type CreateAccountUseCase struct {
	Uow             uow.UowInterface
	EventDispatcher events.EventDispatcherInterface
	AccountCreated  events.EventInterface
}

func NewCreateAccountUseCase(
	Uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	accountCreated events.EventInterface,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		Uow:             Uow,
		EventDispatcher: eventDispatcher,
		AccountCreated:  accountCreated,
	}
}

func (uc *CreateAccountUseCase) Execute(ctx context.Context, input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	output := &CreateAccountOutputDTO{}
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		clientRepository := uc.getClientRepository(ctx)
		accountRepository := uc.getAccountRepository(ctx)

		client, err := clientRepository.FindByID(input.ClientID)
		if err != nil {
			return err
		}

		account, err := entity.NewAccount(client, input.AccountType)
		if err != nil {
			return err
		}

		err = accountRepository.Save(account)
		if err != nil {
			return err
		}

		output.ID = account.ID
		output.ClientID = account.ClientID
		output.AccountType = account.AccountType
		output.CreatedAt = account.CreatedAt
		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.AccountCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.AccountCreated)

	return output, nil
}

func (uc *CreateAccountUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateAccountUseCase) getClientRepository(ctx context.Context) gateway.ClientGateway {
	repo, err := uc.Uow.GetRepository(ctx, "ClientDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.ClientGateway)
}
