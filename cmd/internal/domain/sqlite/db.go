package sqlite

import (
	"4shure/cmd/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	dbPath := filepath.Join("/data", "database.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Appointment{})
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
