package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewClient(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, err := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "John Doe", client.Name)
	assert.Equal(t, "j@j.com", client.Email)
	assert.Equal(t, "bankID", client.BankID)
	assert.Equal(t, clientDateOfBirth, client.DateOfBirth)
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	client, err := NewClient("", "", "", time.Time{})

	assert.NotNil(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	newClientDateOfBirth := time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC)

	err := client.Update("John Doe Update", "j@j.com", newClientDateOfBirth)

	assert.Nil(t, err)
	assert.Equal(t, "John Doe Update", client.Name)
	assert.Equal(t, "j@j.com", client.Email)
	assert.Equal(t, newClientDateOfBirth, client.DateOfBirth)
}

func TestUpdateClientWithInvalidArgs(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	err := client.Update("", "j@j.com", clientDateOfBirth)
	assert.Error(t, err, "Name is required")
}

func TestAddAccountToClient(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	account, err := NewAccount(client, CurrentAccount)
	if err != nil {
		t.Fatalf("Error creating account: %v", err)
	}
	err = client.AddAccount(account)
	if err != nil {
		t.Fatalf("Error adding account to client: %v", err)
	}

	if len(client.Accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(client.Accounts))
	}
}

func TestAddInvalidAccountToClient(t *testing.T) {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := NewClient("John Doe", "j@j.com", "bankID", clientDateOfBirth)

	otherClient, _ := NewClient("Jane Doe", "jane@j.com", "otherBankID", clientDateOfBirth)

	account, _ := NewAccount(otherClient, CurrentAccount)

	err := client.AddAccount(account)
	assert.Error(t, err, "Account does not belong to client")
	assert.Equal(t, 0, len(client.Accounts))
}
