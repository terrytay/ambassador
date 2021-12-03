package main

import (
	"math/rand"

	"github.com/bxcodec/faker/v3"
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/models"
)

func main() {
	database.Connect()

	for i := 0; i < 30; i++ {
		product := models.Product{
			Title:       faker.Username(),
			Description: faker.Username(),
			Image:       faker.URL(),
			Price:       rand.Float64() * 100,
		}

		database.DB.Create(&product)
	}
}
