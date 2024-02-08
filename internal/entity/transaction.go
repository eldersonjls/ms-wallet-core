package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Deposit  TransactionType = "DEPOSIT"
	Removal  TransactionType = "REMOVAL"
	Transfer TransactionType = "TRANSFER"
)

type Transaction struct {
	ID              string
	TransactionType TransactionType
	AccountFrom     *Account `gorm:"foreignKey:AccountFromID"`
	AccountFromID   string
	AccountTo       *Account `gorm:"foreignKey:AccountToID"`
	AccountToID     string
	Amount          float64
	CreatedAt       time.Time
}

func NewTransaction(accountFrom *Account, accountTo *Account, amount float64, transactionType TransactionType) (*Transaction, error) {
	switch transactionType {
	case Deposit:
		if accountTo == nil {
			return nil, errors.New("Destination account is required for deposit")
		}
	case Removal:
		if accountFrom == nil {
			return nil, errors.New("Source account is required for removal")
		}
	case Transfer:
		if accountFrom == nil || accountTo == nil {
			return nil, errors.New("Both source and destination accounts are required for transfer")
		}
	}

	if amount < 0 {
		return nil, errors.New("Amount must be greater than zero")
	}

	switch transactionType {
	case Deposit, Removal, Transfer:
		transaction := &Transaction{
			ID:              uuid.New().String(),
			TransactionType: transactionType,
			AccountFrom:     accountFrom,
			AccountFromID:   accountFrom.ID,
			AccountTo:       accountTo,
			AccountToID:     accountTo.ID,
			Amount:          amount,
			CreatedAt:       time.Now(),
		}

		err := transaction.Validate()
		if err != nil {
			return nil, err
		}

		transaction.Commit()
		return transaction, nil
	default:
		return nil, errors.New("Invalid transaction type")
	}
}

func (t *Transaction) Commit() error {
	if t == nil {
		return errors.New("Cannot commit nil transaction")
	}

	switch t.TransactionType {
	case Deposit:
		if t.AccountTo == nil {
			return errors.New("Destination account is required for deposit transaction")
		}
		t.AccountTo.Credit(t.Amount)
	case Removal:
		if t.AccountFrom == nil {
			return errors.New("Source account is required for removal transaction")
		}
		if t.AccountFrom.Balance < t.Amount {
			return errors.New("Insufficient balance for removal transaction")
		}
		t.AccountFrom.Debit(t.Amount)
	case Transfer:
		if t.AccountFrom == nil || t.AccountTo == nil {
			return errors.New("Both source and destination accounts are required for transfer transaction")
		}
		if t.AccountFrom.Balance < t.Amount {
			return errors.New("Insufficient balance for transfer transaction")
		}
		t.AccountFrom.Debit(t.Amount)
		t.AccountTo.Credit(t.Amount)
	}
	return nil
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("The value must be greater than zero")
	}
	if t.TransactionType == Transfer && t.AccountFrom.Balance < t.Amount {
		return errors.New("Insufficient funds")
	}
	return nil
}
