package postgres

import (
	"log"
	"strings"
	"time"

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
	return "", nil
}

func (m *manager) WriteAuthToken(accountID string, token string) error {
	return nil
}

// Account Stater

func (m *manager) AllAccounts(first int, after int) ([]state.Account, error) {
	return nil, nil
}

func (m *manager) ReadAccount(accountID string) (*state.Account, error) {
	return nil, nil
}

func (m *manager) ReadAccountForEmail(email string) (*state.Account, error) {
	return nil, nil
}

func (m *manager) WriteAccount(in *state.Account, password string) (*state.Account, error) {
	return nil, nil
}

func (m *manager) DeleteAccount(accountID string) error {
	return nil
}

func (m *manager) HistoryForAccount(accountID string) ([]state.Account, error) {
	return nil, nil
}

func (m *manager) RestoreAccount(accountID string, at time.Time) (*state.Account, error) {
	return nil, nil
}
