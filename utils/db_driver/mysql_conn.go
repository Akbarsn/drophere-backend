package db_driver

import (
	"errors"
	"fmt"
	"time"

	"github.com/bccfilkom/drophere-go/utils/env_driver"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func NewMysqlConn(env *env_driver.DatabaseEnvironment) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		env.Username,
		env.Password,
		env.Host,
		env.Port,
		env.Name)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if db.Raw("SELECT 1").RowsAffected < 0 {
		return nil, errors.New("Failed to run test query")
	} else {
		fmt.Println("DB Connected")
	}

	db = db.
		Set("gorm:table_options", "DEFAULT CHARACTER SET=utf8mb4 COLLATE=utf8mb4_general_ci ENGINE=InnoDB").
		Set("gorm:auto_preload", false)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(time.Minute)

	return db, nil
}
