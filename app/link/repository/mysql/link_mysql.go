package mysql

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(conn *gorm.DB) domain.LinkRepository {
	return &UserRepository{DB: conn}
}

func (u UserRepository) Create(l *domain.Link) (*domain.Link, error) {
	panic("implement me")
}

func (u UserRepository) Delete(l *domain.Link) error {
	panic("implement me")
}

func (u UserRepository) FindByID(id uint) (*domain.Link, error) {
	panic("implement me")
}

func (u UserRepository) FindBySlug(slug string) (*domain.Link, error) {
	panic("implement me")
}

func (u UserRepository) ListByUser(userID uint) ([]domain.Link, error) {
	panic("implement me")
}

func (u UserRepository) Update(l *domain.Link) (*domain.Link, error) {
	panic("implement me")
}
