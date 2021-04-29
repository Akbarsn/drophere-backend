package inmemory

import "github.com/bccfilkom/drophere-go/domain"

type UserStorageCredentialRepository struct {
	DB *DB
}

// NewUserStorageCredentialRepository func
func NewUserStorageCredentialRepository(db *DB) domain.UserStorageCredentialRepository {
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
func (repo *UserStorageCredentialRepository) Find(filters domain.UserStorageCredentialFilters, withUserRelation bool) ([]domain.UserStorageCredential, error) {
	creds := make([]domain.UserStorageCredential, 0)
	usersCache := make(map[uint]domain.User)

	// load users first
	if withUserRelation && len(filters.UserIDs) > 0 {
		for _, u := range repo.DB.Users {
			if isInUintSlice(u.ID, filters.UserIDs) {
				usersCache[u.ID] = u
			}
		}
	}

	for _, usc := range repo.DB.UserStorageCredentials {
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
func (repo *UserStorageCredentialRepository) FindByID(id uint, withUserRelation bool) (domain.UserStorageCredential, error) {
	cred := domain.UserStorageCredential{}
	found := false
	for _, usc := range repo.DB.UserStorageCredentials {
		if usc.ID == id {
			cred = usc
			found = true
			break
		}
	}

	if found {
		if withUserRelation {
			for _, u := range repo.DB.Users {
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
func (repo *UserStorageCredentialRepository) Create(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {
	repo.DB.UserStorageCredentials = append(repo.DB.UserStorageCredentials, cred)
	return cred, nil
}

// Update impl
func (repo *UserStorageCredentialRepository) Update(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {

	for i := range repo.DB.UserStorageCredentials {
		if repo.DB.UserStorageCredentials[i].ID == cred.ID {
			repo.DB.UserStorageCredentials[i] = cred
			break
		}
	}

	return cred, nil
}

// Delete impl
func (repo *UserStorageCredentialRepository) Delete(cred domain.UserStorageCredential) error {

	for i := range repo.DB.UserStorageCredentials {
		if repo.DB.UserStorageCredentials[i].ID == cred.ID {
			repo.DB.UserStorageCredentials = append(repo.DB.UserStorageCredentials[:i], repo.DB.UserStorageCredentials[i+1:]...)
			break
		}
	}

	return nil
}
