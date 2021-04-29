package inmemory

import (
	"time"

	"github.com/bccfilkom/drophere-go/domain"
)

// DB struct
type DB struct {
	Users                  []domain.User
	Links                  []domain.Link
	UserStorageCredentials []domain.UserStorageCredential
}

// New func
func New() *DB {
	db := &DB{}
	db.populate()
	return db
}

func str2ptr(s string) *string {
	return &s
}

func time2ptr(t time.Time) *time.Time {
	return &t
}

func (db *DB) populate() {
	db.Users = []domain.User{
		{ID: 1, Email: "user@drophere.link", Name: "User", Password: "123456", DropboxToken: nil, DriveToken: nil},
		{ID: 357, Email: "user_357@drophere.link", Name: "User 357", Password: "123456", DropboxToken: nil, DriveToken: nil},
		{
			ID:                         6631,
			Email:                      "reset+pwd+expired_token@drophere.link",
			Name:                       "Token is set but expired",
			Password:                   "123456",
			RecoverPasswordToken:       str2ptr("expired_recover_password_token"),
			RecoverPasswordTokenExpiry: time2ptr(time.Now().Add(time.Minute * -30)),
		},
		{
			ID:                         12368,
			Email:                      "reset+pwd@drophere.link",
			Name:                       "Token is set",
			Password:                   "123456",
			RecoverPasswordToken:       str2ptr("recover_password_token"),
			RecoverPasswordTokenExpiry: time2ptr(time.Now().Add(time.Minute * 30)),
		},
	}

	db.Links = []domain.Link{
		{ID: 1, UserID: 1, User: &db.Users[0], Title: "Drop file here", Slug: "drop-here", Password: "123098", Description: "drop a file here"},
		{ID: 2, UserID: 1, User: &db.Users[0], Title: "Test Link 2", Slug: "test-link-2", Password: "", Description: "no description"},
		{ID: 3, UserID: 357, User: &db.Users[1], Title: "Another link", Slug: "another-link", Password: "999", Description: "nil here"},
	}

	db.UserStorageCredentials = []domain.UserStorageCredential{
		{
			ID:                 2000,
			UserID:             1,
			ProviderID:         1,
			ProviderCredential: "user_1_mock_token",
			Email:              "user@drophere.link",
			Photo:              "http://my.photo/user1.jpg",
		},
	}
}

// FindUserByEmail func
func (db *DB) FindUserByEmail(email string) (*domain.User, error) {
	for i, u := range db.Users {
		if u.Email == email {
			return &db.Users[i], nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// FindUserByID func
func (db *DB) FindUserByID(id uint) (*domain.User, error) {
	for i, u := range db.Users {
		if u.ID == id {
			return &db.Users[i], nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// CreateUser func
func (db *DB) CreateUser(u *domain.User) (*domain.User, error) {
	db.Users = append(db.Users, *u)
	return u, nil
}
