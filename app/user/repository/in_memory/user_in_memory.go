package in_memory

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
)

type UserRepository struct {
	DB *inmemory.DB
}

func NewUserRepository(db *inmemory.DB) domain.UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) Create(u *domain.User) (*domain.User, error) {
	return ur.DB.CreateUser(u)
}

func (ur *UserRepository) FindByEmail(email string) (*domain.User, error) {
	return ur.DB.FindUserByEmail(email)
}

func (ur *UserRepository) FindByID(id uint) (*domain.User, error) {
	return ur.DB.FindUserByID(id)
}

func (ur *UserRepository) Update(u *domain.User) (*domain.User, error) {
	updated := false
	for i := range ur.DB.Users {
		if ur.DB.Users[i].ID == u.ID {
			ur.DB.Users[i] = *u
			updated = true
			break
		}
	}

	if !updated {
		ur.DB.Users = append(ur.DB.Users, *u)
	}
	return u, nil
}
