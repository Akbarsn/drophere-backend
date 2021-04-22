package mysql

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/jinzhu/gorm"
)

type LinkRepository struct {
	DB *gorm.DB
}

func NewLinkRepository(db *gorm.DB) domain.LinkRepository {
	return &LinkRepository{db}
}

func (lr LinkRepository) Create(l *domain.Link) (*domain.Link, error) {
	if err := lr.DB.Create(l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

func (lr LinkRepository) Delete(l *domain.Link) error {
	return lr.DB.Delete(l).Error
}

func (lr LinkRepository) FindByID(id uint) (*domain.Link, error) {
	l := domain.Link{}
	q := lr.DB.Preload("User").
		Preload("UserStorageCredential").
		Find(&l, id)
	if q.RecordNotFound() {
		return nil, domain.ErrLinkNotFound
	}
	if q.Error != nil {
		return nil, q.Error
	}
	return &l, nil
}

func (lr LinkRepository) FindBySlug(slug string) (*domain.Link, error) {
	l := domain.Link{}
	q := lr.DB.Where("`slug` = ? ", slug).
		Preload("User").
		Preload("UserStorageCredential").
		Find(&l)
	if q.RecordNotFound() {
		return nil, domain.ErrLinkNotFound
	}
	if q.Error != nil {
		return nil, q.Error
	}
	return &l, nil
}

func (lr LinkRepository) ListByUser(userID uint) ([]domain.Link, error) {
	var links []domain.Link
	if err := lr.DB.Where("`user_id` = ? ", userID).
		Preload("User").
		Preload("UserStorageCredential").
		Find(&links).
		Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (lr LinkRepository) Update(l *domain.Link) (*domain.Link, error) {
	if err := lr.DB.Save(l).Error; err != nil {
		return nil, err
	}
	return l, nil
}
