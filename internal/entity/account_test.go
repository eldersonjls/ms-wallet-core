package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewAccount(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC) // Data de nascimento fictícia para o cliente

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	account, err := NewAccount(client, CurrentAccount)

	assert.Nil(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, client, account.Client)
	assert.Equal(t, CurrentAccount, account.AccountType)
	assert.Equal(t, float64(0), account.Balance)
}

func TestCreateNewAccountWithInvalidClient(t *testing.T) {
	account, err := NewAccount(nil, CurrentAccount)

	assert.NotNil(t, err)
	assert.Nil(t, account)
}

func TestCreditAccount(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	account, _ := NewAccount(client, CurrentAccount)

	err := account.Credit(50.0)
	assert.NoError(t, err)
	assert.Equal(t, float64(50), account.Balance)

	err = account.Credit(-20.0)
	assert.Error(t, err)
	assert.Equal(t, float64(50), account.Balance)
}

func TestDebitAccount(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	account, _ := NewAccount(client, CurrentAccount)

	account.Credit(100.0)
	account.Debit(50.0)
	assert.Equal(t, float64(50), account.Balance)
}

func TestDebitAccountWithInsufficientFunds(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	account, _ := NewAccount(client, CurrentAccount)

	account.Credit(50.0)
	err := account.Debit(100.0)
	assert.Error(t, err)
	assert.Equal(t, float64(50), account.Balance)
}

func TestAccountValidation(t *testing.T) {
	clientDateOfBirth := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	account, _ := NewAccount(client, CurrentAccount)

	client.DateOfBirth = time.Date(2008, 1, 1, 0, 0, 0, 0, time.UTC)
	err := account.Validate()
	assert.Error(t, err, "Client must be at least 18 years old for a current account")
}
