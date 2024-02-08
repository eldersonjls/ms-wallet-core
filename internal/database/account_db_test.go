package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
	client    *entity.Client
}

func (s *AccountDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")

	s.Nil(err)

	s.db = db

	db.Exec("create table clients (id varchar(255), name varchar(255), email varchar(255), bank_id varchar(255), date_of_birth date, created_at date, update_at date)")
	db.Exec("create table accounts (id varchar(255), client_id varchar(255), account_type varchar(255), balance int, created_at date, update_at date)")

	s.accountDB = NewAccountDB(db)

	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	s.client, _ = entity.NewClient("John Due", "j@j.com", "123", clientDateOfBirth)
}

func (s *AccountDBTestSuite) TearDownSuite() {
	defer s.db.Close()

	s.db.Exec("drop table clients")
	s.db.Exec("drop table accounts")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSave() {
	account, err := entity.NewAccount(s.client, entity.CurrentAccount)
	s.Require().NoError(err)

	err = s.accountDB.Save(account)
	s.NoError(err)
}

func (s *AccountDBTestSuite) TestFindByID() {
	s.db.Exec("insert into clients (id, name, email, bank_id, date_of_birth, created_at, update_ate) values (?, ?, ?, ?, ?, ?, ?)", s.client.ID, s.client.Name, s.client.Email, s.client.BankID, s.client.DateOfBirth, s.client.CreatedAt, s.client.UpdatedAt)

	account, err := entity.NewAccount(s.client, entity.CurrentAccount)
	err = s.accountDB.Save(account)

	s.Nil(err)

	accountDB, err := s.accountDB.FindByID(account.ID)
	s.Nil(err)

	s.Equal(account.ID, accountDB.ID)
	s.Equal(account.ClientID, accountDB.ClientID)
	s.Equal(account.AccountType, accountDB.AccountType)
	s.Equal(account.Balance, accountDB.Balance)

	s.Equal(account.Client.ID, accountDB.Client.ID)
	s.Equal(account.Client.Name, accountDB.Client.Name)
	s.Equal(account.Client.Email, accountDB.Client.Email)
	s.Equal(account.Client.BankID, accountDB.Client.BankID)
	s.Equal(account.Client.DateOfBirth, accountDB.Client.DateOfBirth)
}
