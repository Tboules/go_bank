package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank port=5433 sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (*PostgresStore) CreateAccount(*Account) error {
	return nil
}

func (*PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (*PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (*PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
