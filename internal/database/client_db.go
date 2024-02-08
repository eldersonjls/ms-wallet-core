package database

import (
	"database/sql"
	"fmt"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
)

type ClientDB struct {
	DB *sql.DB
}

func NewClientDB(db *sql.DB) *ClientDB {
	return &ClientDB{
		DB: db,
	}
}

func (c *ClientDB) FindByID(id string) (*entity.Client, error) {
	client := &entity.Client{}
	stmt, err := c.DB.Prepare("select id, name, email, bank_id, date_of_birth, created_at from clients where id = ?")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(id)
	if err := row.Scan(
		&client.ID,
		&client.Name,
		&client.Email,
		&client.BankID,
		&client.DateOfBirth,
		&client.CreatedAt); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *ClientDB) Save(client *entity.Client) error {
	stmt, err := c.DB.Prepare("insert into clients (id, name, email, bank_id, date_of_birth, created_at) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		client.ID,
		client.Name,
		client.Email,
		client.BankID,
		client.DateOfBirth,
		client.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
