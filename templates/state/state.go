package state

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrRecordNotFound means the record was not found.
	ErrRecordNotFound = errors.New("Record not found")
	// ErrRecordInvalid means the record is invalid.
	ErrRecordInvalid = errors.New("Record invalid")
)

// Stater is the interface that wraps all Stater interfaces.
type Stater interface {
	AuthStater
	AccountStater
}

type AuthStater interface {
	ReadAuthToken(token string) (string, error)
	WriteAuthToken(accountID string, token string) error
}

// AccountStater is the interface that wraps Account I/O.
type AccountStater interface {
	FetchAccounts(first int, after string) (*Accounts, error)
	ReadAccount(accountID string) (*Account, error)
	ReadAccountForEmail(email string) (*Account, error)
	WriteAccount(in *Account, password string) (*Account, error)
	DeleteAccount(accountID string) error
	HistoryForAccount(accountID string) (Accounts, error)
	RestoreAccount(accountID string, at time.Time) (*Account, error)
}

type Backend func(map[string]string) Stater

func Register(kind string, backend Backend) {
	backends[kind] = backend
}

func NewState(kind string, cfg map[string]string) Stater {
	maker, ok := backends[kind]
	if !ok {
		log.Fatalf("State Error: backend '%s' not registered", kind)
	}
	return maker(cfg)
}

var backends = make(map[string]Backend)

// EncryptPassword returns an encrypted password.
func EncryptPassword(in string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(pass)
}
