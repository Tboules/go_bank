package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
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

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id serial PRIMARY KEY,
    first_name varchar(50),
    last_name varchar(50),
    number serial,
    balance serial,
    created_at timestamp
  )`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) (*Account, error) {
	query := `
    INSERT INTO account (first_name, last_name, number, balance, created_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, first_name, last_name, number, balance, created_at
`
	account := new(Account)

	err := s.db.QueryRow(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt).Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return account, nil
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

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accountSlice := []*Account{}

	for rows.Next() {
		account := new(Account)

		if err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.Number, &account.Balance, &account.CreatedAt); err != nil {
			return accountSlice, err
		}

		accountSlice = append(accountSlice, account)

	}

	return accountSlice, nil

}
