package database

import (
	"b2match/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetUpDB() {
	db, err := gorm.Open(sqlite.Open("b2match.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Company{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.Meeting{})
	db.AutoMigrate(&models.Invite{})

	DB = db
}
