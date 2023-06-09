package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := gorm.Open(sqlite.Open("b2match.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}
}
