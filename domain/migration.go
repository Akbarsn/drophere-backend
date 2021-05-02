package domain

import "context"

type MigrationUseCase interface {
	Migrate(ctx context.Context) (string, error)
}

type MigrationRepository interface {
	Migrate() (string, error)
}
