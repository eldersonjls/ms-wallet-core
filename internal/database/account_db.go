package database

import (
	"database/sql"
	"fmt"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
)

type AccountDB struct {
	DB *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDB {
	return &AccountDB{
		DB: db,
	}
}

func (a *AccountDB) FindByID(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client

	stmt, err := a.DB.Prepare("Select a.id, a.client_id, a.account_type, a.balance, a.created_at, c.id, c.name, c.email, c.bank_id, c.date_of_birth, c.created_at FROM accounts a INNER JOIN clients c ON a.client_id = c.id WHERE a.id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	err = row.Scan(
		&account.ID,
		&account.Client.ID,
		&account.AccountType,
		&account.Balance,
		&account.CreatedAt,
		&client.ID,
		&client.Name,
		&client.Email,
		&client.BankID,
		&client.DateOfBirth,
		&client.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *AccountDB) Save(account *entity.Account) error {
	stmt, err := a.DB.Prepare("insert into accounts (id, client_id, account_type, balance, created_at) values (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		account.ID,
		account.Client.ID,
		account.AccountType,
		account.Balance,
		account.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (a *AccountDB) UpdateBalance(account *entity.Account) error {
	stmt, err := a.DB.Prepare("update accounts set balance = ? where id = ?")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(account.Balance, account.ID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func printSQL(query string) {
	fmt.Println("SQL Executado:", query)
}
