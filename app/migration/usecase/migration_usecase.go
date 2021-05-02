package usecase

import (
	"context"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
)

type MigrationUseCase struct {
	Repository     domain.MigrationRepository
	ContextTimeOut time.Duration
}

func NewMigrationUseCase(repo domain.MigrationRepository) domain.MigrationUseCase {
	return &MigrationUseCase{
		Repository: repo,
	}
}

func (m *MigrationUseCase) Migrate(ctx context.Context) (string, error) {
	ctx, cancelContext := context.WithTimeout(ctx, m.ContextTimeOut)
	defer cancelContext()

	result, err := m.Repository.Migrate()
	if err != nil {
		return "", err
	}

	return result, nil
}
