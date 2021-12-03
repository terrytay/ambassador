package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/helpers"
	"github.com/terrytay/ambassador/src/models"
)

func Products(c echo.Context) error {
	var products []models.Product

	database.DB.Find(&products)

	return c.JSON(http.StatusOK, products)
}

func CreateProduct(c echo.Context) error {
	var product models.Product

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Error{Message: "invalid input"})
	}

	database.DB.Create(&product)

	return c.JSON(http.StatusOK, product)
}

func GetProduct(c echo.Context) error {
	var product models.Product

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.Error{Message: "product not found"})
	}

	database.DB.Where("id = ?", id).First(&product)

	if product.Id == 0 {
		return c.JSON(http.StatusNotFound, helpers.Error{Message: "product not found"})
	}

	return c.JSON(http.StatusOK, product)
}
