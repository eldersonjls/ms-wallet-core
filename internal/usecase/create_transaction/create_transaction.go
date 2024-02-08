package create_transaction

import (
	"context"
	"errors"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
	"github.com/eldersonjls/ms-wallet-core/internal/gateway"
	"github.com/eldersonjls/ms-wallet-core/pkg/events"
	"github.com/eldersonjls/ms-wallet-core/pkg/uow"
)

type CreateTransactionInputDTO struct {
	TransactionType entity.TransactionType `json:"transaction_type"`
	AccountIDFrom   string                 `json:"account_id_from"`
	AccountIDTo     string                 `json:"account_id_to"`
	Amount          float64                `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID              string                 `json:"id"`
	TransactionType entity.TransactionType `json:"transaction_type"`
	AccountIDFrom   string                 `json:"account_id_from"`
	AccountIDTo     string                 `json:"account_id_to"`
	Amount          float64                `json:"amount"`
}

type BalanceUpdatedOutputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

type CreateTransactionUseCase struct {
	Uow                uow.UowInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated     events.EventInterface
}

func NewCreateTransactionUseCase(
	Uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow:                Uow,
		EventDispatcher:    eventDispatcher,
		TransactionCreated: transactionCreated,
		BalanceUpdated:     balanceUpdated,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input CreateTransactionInputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}
	balanceUpdatedOutput := &BalanceUpdatedOutputDTO{}
	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountRepository := uc.getAccountRepository(ctx)
		transactionRepository := uc.getTransactionRepository(ctx)

		var accountFrom *entity.Account
		var accountTo *entity.Account
		var err error

		switch input.TransactionType {
		case entity.Deposit:
			if input.AccountIDTo == "" {
				return errors.New("Destination account is required for deposit")
			}
			accountTo, err = accountRepository.FindByID(input.AccountIDTo)
			accountFrom = accountTo
			if err != nil {
				return err
			}
		case entity.Removal:
			if input.AccountIDFrom == "" {
				return errors.New("Source account is required for removal")
			}
			accountFrom, err = accountRepository.FindByID(input.AccountIDFrom)
			accountTo = accountFrom
			if err != nil {
				return err
			}
		case entity.Transfer:
			if input.AccountIDFrom == "" || input.AccountIDTo == "" {
				return errors.New("Both source and destination accounts are required for transfer")
			}
			accountFrom, err = accountRepository.FindByID(input.AccountIDFrom)
			if err != nil {
				return err
			}
			accountTo, err = accountRepository.FindByID(input.AccountIDTo)
			if err != nil {
				return err
			}
		default:
			return errors.New("Invalid transaction type")
		}

		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount, input.TransactionType)
		if err != nil {
			return err
		}

		if accountFrom != nil {
			err = accountRepository.UpdateBalance(accountFrom)
			if err != nil {
				return err
			}
		}

		if accountTo != nil {
			err = accountRepository.UpdateBalance(accountTo)
			if err != nil {
				return err
			}
		}

		err = transactionRepository.Create(transaction)
		if err != nil {
			return err
		}

		output.ID = transaction.ID
		output.TransactionType = input.TransactionType
		output.AccountIDFrom = input.AccountIDFrom
		output.AccountIDTo = input.AccountIDTo
		output.Amount = input.Amount

		if accountFrom != nil && accountTo != nil {
			println("x")
			balanceUpdatedOutput.AccountIDFrom = input.AccountIDFrom
			balanceUpdatedOutput.AccountIDTo = input.AccountIDTo
			balanceUpdatedOutput.BalanceAccountIDFrom = accountFrom.Balance
			balanceUpdatedOutput.BalanceAccountIDTo = accountTo.Balance
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	uc.TransactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.TransactionCreated)

	uc.BalanceUpdated.SetPayload(balanceUpdatedOutput)
	uc.EventDispatcher.Dispatch(uc.BalanceUpdated)
	return output, nil
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (uc *CreateTransactionUseCase) getTransactionRepository(ctx context.Context) gateway.TransactionGateway {
	repo, err := uc.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGateway)
}
