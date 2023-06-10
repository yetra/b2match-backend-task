package main

import (
	"b2match/backend/models"
)

func main() {
	_, err := models.SetUpDatabase()

	if err != nil {
		panic("failed to connect to database")
	}
}
