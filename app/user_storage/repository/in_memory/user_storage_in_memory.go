package in_memory

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
)

type UserStorageCredentialRepository struct {
	DB *inmemory.DB
}

// NewUserStorageCredentialRepository func
func NewUserStorageCredentialRepository(db *inmemory.DB) domain.UserStorageCredentialRepository {
	return &UserStorageCredentialRepository{db}
}

func isInUintSlice(u uint, slice []uint) bool {
	for _, el := range slice {
		if el == u {
			return true
		}
	}
	return false
}

// Find impl
func (usr *UserStorageCredentialRepository) Find(filters domain.UserStorageCredentialFilters, withUserRelation bool) ([]domain.UserStorageCredential, error) {
	creds := make([]domain.UserStorageCredential, 0)
	usersCache := make(map[uint]domain.User)

	// load users first
	if withUserRelation && len(filters.UserIDs) > 0 {
		for _, u := range usr.DB.Users {
			if isInUintSlice(u.ID, filters.UserIDs) {
				usersCache[u.ID] = u
			}
		}
	}

	for _, usc := range usr.DB.UserStorageCredentials {
		if filters.UserIDs != nil && (len(filters.UserIDs) == 0 ||
			!isInUintSlice(usc.UserID, filters.UserIDs)) {
			continue
		}

		if filters.ProviderIDs != nil && (len(filters.ProviderIDs) == 0 ||
			!isInUintSlice(usc.ProviderID, filters.ProviderIDs)) {
			continue
		}

		if withUserRelation {
			usc.User = usersCache[usc.UserID]
		}

		creds = append(creds, usc)
	}

	return creds, nil
}

// FindByID impl
func (usr *UserStorageCredentialRepository) FindByID(id uint, withUserRelation bool) (domain.UserStorageCredential, error) {
	cred := domain.UserStorageCredential{}
	found := false
	for _, usc := range usr.DB.UserStorageCredentials {
		if usc.ID == id {
			cred = usc
			found = true
			break
		}
	}

	if found {
		if withUserRelation {
			for _, u := range usr.DB.Users {
				if u.ID == cred.UserID {
					cred.User = u
					break
				}
			}
		}

		return cred, nil
	}
	return cred, domain.ErrUserStorageCredentialNotFound
}

// Create impl
func (usr *UserStorageCredentialRepository) Create(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {
	usr.DB.UserStorageCredentials = append(usr.DB.UserStorageCredentials, cred)
	return cred, nil
}

// Update impl
func (usr *UserStorageCredentialRepository) Update(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {

	for i := range usr.DB.UserStorageCredentials {
		if usr.DB.UserStorageCredentials[i].ID == cred.ID {
			usr.DB.UserStorageCredentials[i] = cred
			break
		}
	}

	return cred, nil
}

// Delete impl
func (usr *UserStorageCredentialRepository) Delete(cred domain.UserStorageCredential) error {

	for i := range usr.DB.UserStorageCredentials {
		if usr.DB.UserStorageCredentials[i].ID == cred.ID {
			usr.DB.UserStorageCredentials = append(usr.DB.UserStorageCredentials[:i], usr.DB.UserStorageCredentials[i+1:]...)
			break
		}
	}

	return nil
}
