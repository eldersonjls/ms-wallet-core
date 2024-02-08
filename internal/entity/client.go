package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID          string
	Name        string
	Email       string
	BankID      string
	DateOfBirth time.Time
	Accounts    []*Account
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewClient(name string, email string, bankID string, dateOfBirth time.Time) (*Client, error) {
	client := &Client{
		ID:          uuid.New().String(),
		Name:        name,
		Email:       email,
		BankID:      bankID,
		DateOfBirth: dateOfBirth,
		CreatedAt:   time.Now(),
	}
	err := client.Validate()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) Validate() error {
	if c.Name == "" {
		return errors.New("Name is required")
	}
	if c.Email == "" {
		return errors.New("Email is required")
	}
	if c.BankID == "" {
		return errors.New("Bank is required")
	}
	return nil
}

func (c *Client) Update(name string, email string, dateOfBirth time.Time) error {
	c.Name = name
	c.Email = email
	c.DateOfBirth = dateOfBirth
	c.UpdatedAt = time.Now()
	err := c.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) AddAccount(account *Account) error {
	if account == nil {
		return errors.New("Account is nil")
	}
	if account.Client.ID != c.ID {
		return errors.New("Account does not belong to client")
	}
	c.Accounts = append(c.Accounts, account)
	return nil
}
