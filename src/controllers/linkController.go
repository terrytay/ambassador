package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/helpers"
	"github.com/terrytay/ambassador/src/models"
)

func Link(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.GenericResponse{Message: "invalid id"})
	}

	var links []models.Link

	database.DB.Where("user_id = ?", id).Find(&links)

	return c.JSON(http.StatusOK, links)
}
