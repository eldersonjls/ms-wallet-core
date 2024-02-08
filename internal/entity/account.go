package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	CurrentAccount AccountType = "CURRENT"
	SalaryAccount  AccountType = "SALARY"
	SavingsAccount AccountType = "SAVINGS"
)

type Account struct {
	ID           string
	Client       *Client
	ClientID     string
	AccountType  AccountType
	Balance      float64
	Transactions []Transaction
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (a *Account) Validate() error {
	if a.ID == "" {
		return errors.New("Account ID is required")
	}
	if a.Client == nil || a.Client.ID == "" {
		return errors.New("A valid client is required for the account")
	}
	if a.AccountType == "" {
		return errors.New("Account type is required")
	}
	if a.Balance < 0 {
		return errors.New("Account balance cannot be negative")
	}

	switch a.AccountType {
	case SalaryAccount:
		// Por exemplo, contas salário não podem ter saldo negativo
		if a.Balance < 0 {
			return errors.New("Salary account cannot have a negative balance")
		}
	case SavingsAccount:
		// Exemplo: contas poupança podem exigir um depósito mínimo inicial
		const minInitialBalance = 100.0
		if a.Balance < minInitialBalance {
			return errors.New("savings account must have a minimum initial balance")
		}
	case CurrentAccount:
		// Exemplo: verificação de idade mínima do cliente para contas correntes
		if !isValidAge(a.Client.DateOfBirth) {
			return errors.New("Client must be at least 18 years old for a current account")
		}
	}
	return nil
}

func isValidAge(dob time.Time) bool {
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}
	return age >= 18 // Verifica se a idade é pelo menos 18 anos
}

func NewAccount(client *Client, accountType AccountType) (*Account, error) {
	if client == nil {
		return nil, errors.New("Client is required")
	}

	err := client.Validate()
	if err != nil {
		return nil, err
	}

	switch accountType {
	case CurrentAccount, SalaryAccount, SavingsAccount:
		account := &Account{
			ID:           uuid.New().String(),
			Client:       client,
			AccountType:  accountType,
			Balance:      0.0,
			Transactions: make([]Transaction, 0),
			CreatedAt:    time.Now(),
		}

		err = account.Validate()
		if err != nil {
			return nil, err
		}

		return account, nil
	default:
		return nil, errors.New("Invalid account type")
	}
}

func (a *Account) Credit(amount float64) error {
	if amount < 0 {
		return errors.New("Amount must be positive")
	}
	a.Balance += amount
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Account) Debit(amount float64) error {
	if amount < 0 {
		return errors.New("Amount must be positive")
	}
	if a.Balance < amount {
		return errors.New("Insufficient funds")
	}
	a.Balance -= amount
	a.UpdatedAt = time.Now()
	return nil
}
