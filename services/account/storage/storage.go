package storage

import (
	"time"
	"gopkg.in/go-playground/validator.v9"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	validate = validator.New()

	ErrAccountNotFound = errors.New("account not found")
	ErrEmailExists     = errors.New("email already exists")
	ErrNoDatabase      = errors.New("no database connection details")
	ErrNoPasswordGiven = errors.New("a password is required")
)

type (
	Account struct {
		ID                 string    `db:"id"`
		Name               string    `validate:"required"`
		Email              string    `validate:"required"`
		HashedPassword     string    `db:"hashed_password"`
		ConfirmationToken  string
		PasswordResetToken string
		Metadata           map[string]string
		CreatedAt          time.Time `db:"created_at"`
	}

	Storage interface {
		List(count int32, token string) ([]*Account, string, error)
		ReadByID(ID string) (*Account, error)
		ReadByEmail(email string) (*Account, error)
		Create(a *Account, password string) error
		Update(a *Account) error
		Delete(ID string) error
		Confirm(token string) (*Account, error)
		GeneratePasswordToken(email string) (*Account, error)
		UpdatePassword(string, string) (*Account, error)
		Migrate() error
		Truncate() error
		Close() error
	}
)

func (a *Account) Valid() error {
	return validate.Struct(a)
}

func (a *Account) HashPassword(password string) error {
	if password == "" {
		return ErrNoPasswordGiven
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	a.HashedPassword = string(hash[:])

	return nil
}

func (a *Account) ComparePasswordToHash(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.HashedPassword), []byte(password))
}
