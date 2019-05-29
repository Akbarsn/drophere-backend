package domain

import (
	"errors"
	"time"
)

var (
	// ErrUserInvalidPassword error
	ErrUserInvalidPassword = errors.New("user: invalid password")
	// ErrUserNotFound error
	ErrUserNotFound = errors.New("user: not found")
)

// User model
type User struct {
	ID           uint
	Email        string
	Name         string
	Password     string
	DropboxToken *string
	DriveToken   *string
}

// SetPassword hash input password and set it to the user struct
func (u *User) SetPassword(password string, hasher Hasher) {
	var err error
	u.Password, err = hasher.Hash(password)
	if err != nil {
		panic(err)
	}
}

// VerifyPassword checks if the encrypted password content is
// equal to the given plain password
func (u *User) VerifyPassword(plainPwd string, hasher Hasher) bool {
	return hasher.Verify(u.Password, plainPwd)
}

// UserCredentials model
type UserCredentials struct {
	Token  string
	Expiry *time.Time
}

// UserService abstraction
type UserService interface {
	Register(email, name, password string) (*User, error)
	Auth(email, password string) (*UserCredentials, error)
	Update(userID uint, name, password, oldPassword *string) (*User, error)
}

// UserRepository abstraction
type UserRepository interface {
	Create(u *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(u *User) (*User, error)
}

// Authenticator is external authentication service
type Authenticator interface {
	Authenticate(u *User) (*UserCredentials, error)
}
