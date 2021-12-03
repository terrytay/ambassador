package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/models"
)

func Ambassador(c echo.Context) error {
	var users []models.User

	database.DB.Where("is_ambassador = true").Find(&users)

	return c.JSON(http.StatusOK, users)
}