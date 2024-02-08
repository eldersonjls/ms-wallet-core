package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	client        *entity.Client
	client2       *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
	transactionDB *TransactionDB
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)

	s.db = db

	db.Exec("create table clients (id varchar(255), name varchar(255), email varchar(255), bank_id varchar(255), date_of_birth date, created_at date, update_at date)")
	db.Exec("create table accounts (id varchar(255), client_id varchar(255), account_type varchar(255), balance int, created_at date, update_at date)")
	db.Exec("create table transactions (id varchar(255), transaction_type varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount int, created_at date)")

	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, err := entity.NewClient("John Due", "j@j.com", "123", clientDateOfBirth)
	s.Nil(err)
	s.client = client

	client2, err := entity.NewClient("Jeny Due", "jj@j.com", "1234", clientDateOfBirth)
	s.Nil(err)
	s.client2 = client2

	accountFrom, err := entity.NewAccount(s.client, entity.CurrentAccount)
	accountFrom.Balance = 100
	s.accountFrom = accountFrom

	accountTo, err := entity.NewAccount(s.client2, entity.CurrentAccount)
	accountTo.Balance = 100
	s.accountTo = accountTo

	s.transactionDB = NewTransactionDB(db)
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("drop table clients")
	s.db.Exec("drop table accounts")
	s.db.Exec("drop table transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.accountFrom, s.accountTo, 100, entity.Deposit)
	s.Nil(err)

	err = s.transactionDB.Create(transaction)
	s.Nil(err)
}
