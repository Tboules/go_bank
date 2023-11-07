package main

import (
	"database/sql"
	"fmt"

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

	account, err := extractAccountFromRow(
		s.db.QueryRow(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt),
	)

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	rows, err := s.db.Exec(`DELETE FROM account WHERE id = $1`, id)

	count, err := rows.RowsAffected()

	if err != nil || count == 0 {

		return fmt.Errorf("Account with id:%v not found", id)

	}

	return err
}

func (*PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	query := `
    SELECT * FROM account
    WHERE id = $1
  `

	account, err := extractAccountFromRow(
		s.db.QueryRow(query, id),
	)

	if err != nil {
		return nil, fmt.Errorf("Account %d not found", id)
	}

	return account, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accountSlice := []*Account{}

	for rows.Next() {
		account, err := extractAccountFromRow(rows)
		if err != nil {
			return nil, err
		}

		accountSlice = append(accountSlice, account)
	}

	return accountSlice, nil

}

type SqlRowOrRows interface {
	Scan(dest ...any) error
}

func extractAccountFromRow(rowOrRows SqlRowOrRows) (*Account, error) {
	account := new(Account)

	err := rowOrRows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
