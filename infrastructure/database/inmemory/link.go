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
func (lr *LinkRepository) Create(l *domain.Link) (*domain.Link, error) {
	l.ID = uint(len(lr.DB.Links) + 1)
	lr.DB.Links = append(lr.DB.Links, *l)
	return l, nil
}

// Delete implementation
func (lr *LinkRepository) Delete(l *domain.Link) error {

	for i := range lr.DB.Links {
		if lr.DB.Links[i].ID == l.ID {
			lr.DB.Links = append(lr.DB.Links[:i], lr.DB.Links[i+1:]...)
			break
		}
	}

	return nil
}

// FindByID implementation
func (lr *LinkRepository) FindByID(id uint) (*domain.Link, error) {
	for i := range lr.DB.Links {
		if lr.DB.Links[i].ID == id {
			return &lr.DB.Links[i], nil
		}
	}

	return nil, domain.ErrLinkNotFound
}

// FindBySlug implementation
func (lr *LinkRepository) FindBySlug(slug string) (*domain.Link, error) {
	for i := range lr.DB.Links {
		if lr.DB.Links[i].Slug == slug {
			return &lr.DB.Links[i], nil
		}
	}

	return nil, domain.ErrLinkNotFound
}

// ListByUser implementation
func (lr *LinkRepository) ListByUser(userID uint) ([]domain.Link, error) {
	links := make([]domain.Link, 0, len(lr.DB.Links))
	for _, link := range lr.DB.Links {
		if link.UserID == userID {
			links = append(links, link)
		}
	}

	return links, nil
}

// Update implementation
func (lr *LinkRepository) Update(l *domain.Link) (link *domain.Link, err error) {
	link = l
	for i := range lr.DB.Links {
		if lr.DB.Links[i].ID == l.ID {
			lr.DB.Links[i] = *l
			return
		}
	}
	lr.DB.Links = append(lr.DB.Links, *l)
	return
}
