package database

import (
	"database/sql"
	"fmt"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
)

type TransactionDB struct {
	DB *sql.DB
}

func NewTransactionDB(db *sql.DB) *TransactionDB {
	return &TransactionDB{
		DB: db,
	}
}

func (t *TransactionDB) Create(transaction *entity.Transaction) error {
	stmt, err := t.DB.Prepare("insert into transactions (id, transaction_type, account_id_from, account_id_to, amount, created_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		transaction.ID,
		transaction.TransactionType,
		transaction.AccountFrom.ID,
		transaction.AccountTo.ID,
		transaction.Amount,
		transaction.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
