package inmemory

import "github.com/bccfilkom/drophere-go/domain"

type LinkRepository struct {
	DB *DB
}

// NewLinkRepository func
func NewLinkRepository(db *DB) domain.LinkRepository {
	return &LinkRepository{db}
}

// Create implementation
func (repo *LinkRepository) Create(l *domain.Link) (*domain.Link, error) {
	l.ID = uint(len(repo.DB.Links) + 1)
	repo.DB.Links = append(repo.DB.Links, *l)
	return l, nil
}

// Delete implementation
func (repo *LinkRepository) Delete(l *domain.Link) error {

	for i := range repo.DB.Links {
		if repo.DB.Links[i].ID == l.ID {
			repo.DB.Links = append(repo.DB.Links[:i], repo.DB.Links[i+1:]...)
			break
		}
	}

	return nil
}

// FindByID implementation
func (repo *LinkRepository) FindByID(id uint) (*domain.Link, error) {
	for i := range repo.DB.Links {
		if repo.DB.Links[i].ID == id {
			return &repo.DB.Links[i], nil
		}
	}

	return nil, domain.ErrLinkNotFound
}

// FindBySlug implementation
func (repo *LinkRepository) FindBySlug(slug string) (*domain.Link, error) {
	for i := range repo.DB.Links {
		if repo.DB.Links[i].Slug == slug {
			return &repo.DB.Links[i], nil
		}
	}

	return nil, domain.ErrLinkNotFound
}

// ListByUser implementation
func (repo *LinkRepository) ListByUser(userID uint) ([]domain.Link, error) {
	links := make([]domain.Link, 0, len(repo.DB.Links))
	for _, link := range repo.DB.Links {
		if link.UserID == userID {
			links = append(links, link)
		}
	}

	return links, nil
}

// Update implementation
func (repo *LinkRepository) Update(l *domain.Link) (link *domain.Link, err error) {
	link = l
	for i := range repo.DB.Links {
		if repo.DB.Links[i].ID == l.ID {
			repo.DB.Links[i] = *l
			return
		}
	}
	repo.DB.Links = append(repo.DB.Links, *l)
	return
}
