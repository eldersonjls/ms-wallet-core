package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	clientFrom, _ := NewClient("John Doe", "from@example.com", "bankID", clientDateOfBirth)
	clientTo, _ := NewClient("Jane Doe", "to@example.com", "bankID", clientDateOfBirth)

	accountFrom, _ := NewAccount(clientFrom, CurrentAccount)
	accountFrom.Balance = 200.0

	accountTo, _ := NewAccount(clientTo, CurrentAccount)
	accountTo.Balance = 100.0

	t.Run("Successful transaction", func(t *testing.T) {
		transaction, err := NewTransaction(accountFrom, accountTo, 100.0, Transfer)

		assert.Nil(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, Transfer, transaction.TransactionType)
		assert.Equal(t, accountFrom.ID, transaction.AccountFrom.ID)
		assert.Equal(t, accountTo.ID, transaction.AccountTo.ID)
	})
}

func TestNewTransaction_ValidAmount(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	clientFrom, _ := NewClient("John Doe", "from@example.com", "bankID", clientDateOfBirth)
	clientTo, _ := NewClient("Jane Doe", "to@example.com", "bankID", clientDateOfBirth)

	accountFrom, _ := NewAccount(clientFrom, CurrentAccount)
	accountFrom.Balance = 100.0

	accountTo, _ := NewAccount(clientTo, CurrentAccount)
	accountTo.Balance = 100.0

	transaction, err := NewTransaction(accountFrom, accountTo, 50.0, Transfer)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(50.0), transaction.Amount)
}

func TestNewTransaction_Deposit(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	accountFrom, _ := NewAccount(client, CurrentAccount)
	accountFrom.Balance = 100.0

	accountTo := &Account{}

	transaction, err := NewTransaction(accountFrom, accountTo, 100.0, Deposit)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(100.0), accountTo.Balance)
	assert.Equal(t, Deposit, transaction.TransactionType)
}

func TestNewTransaction_Removal(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	accountFrom, _ := NewAccount(client, CurrentAccount)
	accountFrom.Balance = 100.0

	accountTo := &Account{}

	transaction, err := NewTransaction(accountFrom, accountTo, 50.0, Removal)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(50.0), accountFrom.Balance)
	assert.Equal(t, Removal, transaction.TransactionType)
}

func TestNewTransaction_Transfer(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	accountFrom, _ := NewAccount(client, CurrentAccount)
	accountFrom.Balance = 100.0

	accountTo := &Account{}

	transaction, err := NewTransaction(accountFrom, accountTo, 50.0, Transfer)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(50.0), accountFrom.Balance)
	assert.Equal(t, float64(50.0), accountTo.Balance)
	assert.Equal(t, Transfer, transaction.TransactionType)
}

func TestTransaction_Commit_Deposit(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	clientTo, _ := NewClient("Jane Doe", "to@example.com", "bankID", clientDateOfBirth)

	accountFrom := &Account{}

	accountTo, _ := NewAccount(clientTo, CurrentAccount)
	accountTo.Balance = 100.0

	transaction, err := NewTransaction(accountFrom, accountTo, 100.0, Deposit)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(200), accountTo.Balance)
}

func TestTransaction_Commit_Removal(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	clientFrom, _ := NewClient("John Doe", "from@example.com", "bankID", clientDateOfBirth)

	accountFrom, _ := NewAccount(clientFrom, CurrentAccount)
	accountFrom.Balance = 50.0

	accountTo := &Account{}

	transaction, err := NewTransaction(accountFrom, accountTo, 50.0, Removal)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(0), accountFrom.Balance)
}

func TestTransaction_Commit_Transfer(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	clientFrom, _ := NewClient("John Doe", "from@example.com", "bankID", clientDateOfBirth)
	clientTo, _ := NewClient("Jane Doe", "to@example.com", "bankID", clientDateOfBirth)

	accountFrom, _ := NewAccount(clientFrom, CurrentAccount)
	accountFrom.Balance = 50.0

	accountTo, _ := NewAccount(clientTo, CurrentAccount)
	accountTo.Balance = 50.0

	transaction, err := NewTransaction(accountFrom, accountTo, 50.0, Transfer)

	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, float64(0), accountFrom.Balance)
	assert.Equal(t, float64(100), accountTo.Balance)
}
