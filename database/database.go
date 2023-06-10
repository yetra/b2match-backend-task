package database

import (
	"b2match/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetUpDB() {
	DB, err := gorm.Open(sqlite.Open("b2match.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	DB.AutoMigrate(&models.Company{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Event{})
	DB.AutoMigrate(&models.Meeting{})
	DB.AutoMigrate(&models.Invite{})
}
