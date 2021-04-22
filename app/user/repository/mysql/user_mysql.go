package mysql

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db}
}

func (ur UserRepository) Create(user *domain.User) (*domain.User, error) {
	if err := ur.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur UserRepository) FindByEmail(email string) (*domain.User, error) {
	user := domain.User{}
	if q := ur.DB.Where("`email` = ? ", email).Find(&user); q.RecordNotFound() {
		return nil, domain.ErrUserNotFound
	} else if q.Error != nil {
		return nil, q.Error
	}
	return &user, nil
}

func (ur UserRepository) FindByID(id uint) (*domain.User, error) {
	user := domain.User{}
	q := ur.DB.Find(&user, id)
	if q.RecordNotFound() {
		return nil, domain.ErrUserNotFound
	}
	if q.Error != nil {
		return nil, q.Error
	}
	return &user, nil
}

func (ur UserRepository) Update(u *domain.User) (*domain.User, error) {
	if err := ur.DB.Save(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}
