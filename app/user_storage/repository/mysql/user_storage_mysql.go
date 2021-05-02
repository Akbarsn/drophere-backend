package mysql

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/jinzhu/gorm"
)

type UserStorageRepository struct {
	DB *gorm.DB
}

func NewUserStorageRepository(db *gorm.DB) domain.UserStorageCredentialRepository {
	return &UserStorageRepository{db}
}

func (ur *UserStorageRepository) Find(filters domain.UserStorageCredentialFilters, withUserRelation bool) ([]domain.UserStorageCredential, error) {
	var (
		credentials []domain.UserStorageCredential
		query       = ur.DB
	)

	if withUserRelation {
		query = query.Preload("User")
	}

	if filters.UserIDs != nil && len(filters.UserIDs) > 0 {
		query = query.Where("`user_id` IN (?)", filters.UserIDs)
	}

	if filters.ProviderIDs != nil && len(filters.ProviderIDs) > 0 {
		query = query.Where("`provider_id` IN (?)", filters.ProviderIDs)
	}

	err := query.Find(&credentials).Error
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

func (ur *UserStorageRepository) FindByID(id uint, withUserRelation bool) (domain.UserStorageCredential, error) {
	var (
		credential domain.UserStorageCredential
		query      = ur.DB
	)

	if withUserRelation {
		query = query.Preload("User")
	}

	if q := query.Find(&credential, id); q.RecordNotFound() {
		return credential, domain.ErrUserStorageCredentialNotFound
	} else if q.Error != nil {
		return credential, q.Error
	}

	return credential, nil
}

func (ur *UserStorageRepository) Create(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {
	err := ur.DB.Create(&cred).Error
	if err != nil {
		return domain.UserStorageCredential{}, err
	}
	return cred, nil
}

func (ur *UserStorageRepository) Update(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {
	err := ur.DB.Save(&cred).Error
	if err != nil {
		return domain.UserStorageCredential{}, err
	}
	return cred, nil
}

func (ur *UserStorageRepository) Delete(cred domain.UserStorageCredential) error {
	return ur.DB.Delete(&cred).Error
}
