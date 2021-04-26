package db_driver

import (
	"fmt"
	"time"

	"github.com/bccfilkom/drophere-go/utils/env_driver"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(env env_driver.DatabaseEnvironment) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		env.Username,
		env.Password,
		env.Host,
		env.Port,
		env.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db = db.
		Set("gorm:table_options", "DEFAULT CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci ENGINE=InnoDB").
		Set("gorm:auto_preload", false)

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute)

	return db, nil
}
