package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chahiim/ticket_tracker/internal/validator"
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

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(validator.NotBlank(user.Name), "name", "must be provided")
	v.Check(validator.MaxLength(user.Name, 50), "name", "must not be more than 50 bytes long")
	v.Check(validator.NotBlank(user.Email), "email", "must be provided")
	v.Check(validator.MaxLength(user.Email, 50), "email", "must not be more than 50 bytes long")
	v.Check(validator.IsValidEmail(user.Email), "email", "must be a valid email")
	v.Check(validator.NotBlank(string(user.HashedPassword)), "password", "must be provided")
}

func (m *UserModel) Insert(user *User) error {

	query := `
		INSERT INTO user (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, activated`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.HashedPassword,
	).Scan(&user.ID, &user.Created, $user.Active)
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var HashedPassword []byte

	query := `
		SELECT id, password_hash
		FROM users
		WHERE email = $1
		AND activated = TRUE`

	err := m.DB.QueryRow(s, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, user.ErrInvalidCredentials
		}else {
			return 0, err
		}
	}
	//check the password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, user.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	return nil, nil
}
