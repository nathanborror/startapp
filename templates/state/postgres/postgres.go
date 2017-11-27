package postgres

import (
	"log"
	"strings"
	"time"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nathanborror/{{.Name}}/pkg/ledger"
	"github.com/nathanborror/{{.Name}}/state"

	_ "github.com/lib/pq" // postgres backend
)

type manager struct {
	db *sqlx.DB
}

// NewState returns a postgres Stater implementation.
func NewState(cfg map[string]string) state.Stater {
	var args []string
	if cfg["User"] != "" {
		args = append(args, "user="+cfg["User"])
	}
	if cfg["Password"] != "" {
		args = append(args, "password="+cfg["Password"])
	}
	if cfg["Database"] == "" {
		log.Fatal("Missing database name")
	}
	args = append(args, "dbname="+cfg["Database"], "sslmode=disable")
	db, err := sqlx.Connect("postgres", strings.Join(args, " "))
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(ledger.Schema)
	db.MustExec(`
		CREATE TABLE IF NOT EXISTS edge (
			added_id serial PRIMARY KEY,
			from_id uuid NOT NULL,
			to_id uuid NOT NULL,
			kind varchar(255) NOT NULL,
			UNIQUE(from_id, to_id, kind)
		);
		CREATE TABLE IF NOT EXISTS index_account_email (
			added_id serial PRIMARY KEY,
			account_id uuid NOT NULL UNIQUE,
			email varchar(255) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS index_account_authtoken (
			added_id serial PRIMARY KEY,
			account_id uuid NOT NULL UNIQUE,
			token varchar(40) NOT NULL
		)`)
	return &manager{db}
}

// Auth Stater

func (m *manager) ReadAuthToken(token string) (string, error) {
	return "", fmt.Errorf("Not Implemented")
}

func (m *manager) WriteAuthToken(accountID string, token string) error {
	return fmt.Errorf("Not Implemented")
}

// Account Stater

func (m *manager) FetchAccounts(first int, after string) (*state.Accounts, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (m *manager) ReadAccount(accountID string) (*state.Account, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (m *manager) ReadAccountForEmail(email string) (*state.Account, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (m *manager) WriteAccount(in *state.Account, password string) (*state.Account, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (m *manager) DeleteAccount(accountID string) error {
	return fmt.Errorf("Not Implemented")
}

func (m *manager) HistoryForAccount(accountID string) (*state.Accounts, error) {
	return nil, fmt.Errorf("Not Implemented")
}

func (m *manager) RestoreAccount(accountID string, at time.Time) (*state.Account, error) {
	return nil, fmt.Errorf("Not Implemented")
}

// Private

func wrapErr(err error) error {
	if err == sql.ErrNoRows {
		return state.ErrRecordNotFound
	}
	if err == ledger.ErrRecordNotFound {
		return state.ErrRecordNotFound
	}
	return err
}
