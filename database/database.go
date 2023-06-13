package database

import (
	"b2match/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Company{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.Meeting{})
	db.AutoMigrate(&models.Invite{})
}

func SetUpDB(dsn string) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect to database")
	}

	migrate(db)

	DB = db
}
