package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/jinzhu/gorm"
)

type MigrationRepository struct {
	DB *gorm.DB
}

func NewMigrationRepository(db *gorm.DB) domain.MigrationRepository {
	return &MigrationRepository{
		DB: db,
	}
}

func (m *MigrationRepository) Migrate() (string, error) {
	wd, err := os.Getwd()
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/files/sql/migrate.sql", wd))
	if err != nil {
		return "", err
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		result := m.DB.Exec(request)

		if result.Error != nil {
			return "", err
		}
	}

	return "Migrate Successfully", nil
}
