package data

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/chahiim/ticket_tracker/internal/validator"
)

var (
	ErrRecordNotFound     = errors.New("user: No Matching Record Found")
	ErrInvalidCredentials = errors.New("user: Invalid Credentials")
	ErrDuplicatedEmail    = errors.New("user: Duplicate Email")
)

// User struct
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}
