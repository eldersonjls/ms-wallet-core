package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/eldersonjls/ms-wallet-core/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	clientDB *ClientDB
}

func (s *ClientDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)

	s.db = db

	db.Exec("create table clients (id varchar(255), name varchar(255), email varchar(255), bank_id varchar(255), date_of_birth date, created_at date, update_at date)")

	s.clientDB = NewClientDB(db)
}

func (s *ClientDBTestSuite) TearDownSuite() {
	defer s.db.Close()

	s.db.Exec("drop table clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (s *ClientDBTestSuite) TestSave() {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client := &entity.Client{
		ID:          "1",
		Name:        "John Due",
		Email:       "j@j.com",
		BankID:      "123",
		DateOfBirth: clientDateOfBirth,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.clientDB.Save(client)

	s.Nil(err)
}

func (s *ClientDBTestSuite) TestGet() {
	clientDateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	client, _ := entity.NewClient("John Due", "j@j.com", "123", clientDateOfBirth)
	s.clientDB.Save(client)

	clientDB, err := s.clientDB.FindByID(client.ID)

	s.Nil(err)
	s.Equal(client.ID, clientDB.ID)
	s.Equal(client.Name, clientDB.Name)
	s.Equal(client.Email, clientDB.Email)
}
