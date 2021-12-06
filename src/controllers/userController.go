package controllers

import (
	"context"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/models"
)

func Ambassador(c echo.Context) error {
	var users []models.User

	database.DB.Where("is_ambassador = true").Find(&users)

	return c.JSON(http.StatusOK, users)
}

func Rankings(c echo.Context) error {
	rankings, err := database.Cache.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err,
		})
	}

	result := make(map[string]float64, len(rankings))

	for _, ranking := range rankings {
		result[ranking.Member.(string)] = ranking.Score
	}

	return c.JSON(http.StatusOK, result)
}
