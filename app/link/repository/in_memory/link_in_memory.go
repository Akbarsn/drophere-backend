package in_memory

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
)

type LinkRepository struct {
	DB *inmemory.DB
}

// NewLinkRepository func
func NewLinkRepository(db *inmemory.DB) domain.LinkRepository {
	return &LinkRepository{db}
}

func (lr *LinkRepository) Create(l *domain.Link) (*domain.Link, error) {
	l.ID = uint(len(lr.DB.Links) + 1)
	lr.DB.Links = append(lr.DB.Links, *l)
	return l, nil
}

func (lr *LinkRepository) Delete(l *domain.Link) error {
	for i := range lr.DB.Links {
		if lr.DB.Links[i].ID == l.ID {
			lr.DB.Links = append(lr.DB.Links[:i], lr.DB.Links[i+1:]...)
			break
		}
	}

	return nil
}

func (lr *LinkRepository) FindByID(id uint) (*domain.Link, error) {
	panic("implement me")
	for i := range lr.DB.Links {
		if lr.DB.Links[i].ID == id {
			return &lr.DB.Links[i], nil
		}
	}

	return nil, domain.ErrLinkNotFound
}

func (lr *LinkRepository) FindBySlug(slug string) (*domain.Link, error) {
	for i := range lr.DB.Links {
		if lr.DB.Links[i].Slug == slug {
			return &lr.DB.Links[i], nil
		}
	}

	return nil, domain.ErrLinkNotFound
}

func (lr *LinkRepository) ListByUser(userID uint) ([]domain.Link, error) {
	links := make([]domain.Link, 0, len(lr.DB.Links))
	for _, link := range lr.DB.Links {
		if link.UserID == userID {
			links = append(links, link)
		}
	}

	return links, nil
}

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
