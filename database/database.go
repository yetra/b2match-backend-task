package database

import (
	"b2match/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Company{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.Meeting{})
	db.AutoMigrate(&models.Invite{})
}

func SetUpDB(dsn string, config *gorm.Config) {
	db, err := gorm.Open(sqlite.Open(dsn), config)

	if err != nil {
		panic("failed to connect to database")
	}

	migrate(db)

	DB = db
}
