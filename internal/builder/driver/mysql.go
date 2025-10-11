package driver

import (
	"fmt"

	"github.com/aburizalpurnama/travel/internal/config"
	_ "github.com/go-sql-driver/mysql" // defines mysql driver used
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewMysqlDatabase return gorm dbmap object with MySQL options param
func NewMysqlDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s", cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBTimezone)), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	mysqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	mysqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	mysqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	mysqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)

	err = mysqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
