package inmemory

import "github.com/bccfilkom/drophere-go/domain"

type userRepository struct {
	db *DB
}

// NewUserRepository func
func NewUserRepository(db *DB) domain.UserRepository {
	return &userRepository{db}
}

// Create implementation
func (repo *userRepository) Create(user *domain.User) (*domain.User, error) {
	return repo.db.CreateUser(user)
}

// FindByEmail implementation
func (repo *userRepository) FindByEmail(email string) (*domain.User, error) {
	return repo.db.FindUserByEmail(email)
}

// FindByID implementation
func (repo *userRepository) FindByID(id uint) (*domain.User, error) {
	return repo.db.FindUserByID(id)
}

// Update implementation
func (repo *userRepository) Update(u *domain.User) (*domain.User, error) {
	updated := false
	for i := range repo.db.Users {
		if repo.db.Users[i].ID == u.ID {
			repo.db.Users[i] = *u
			updated = true
			break
		}
	}

	if !updated {
		repo.db.Users = append(repo.db.Users, *u)
	}
	return u, nil
}
