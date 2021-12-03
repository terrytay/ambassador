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
		return c.JSON(http.StatusBadRequest, helpers.GenericResponse{Message: "invalid input"})
	}

	database.DB.Create(&product)

	return c.JSON(http.StatusOK, product)
}

func GetProduct(c echo.Context) error {
	var product models.Product

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.GenericResponse{Message: "product not found"})
	}

	database.DB.Where("id = ?", id).First(&product)

	if product.Id == 0 {
		return c.JSON(http.StatusNotFound, helpers.GenericResponse{Message: "product not found"})
	}

	return c.JSON(http.StatusOK, product)
}

func UpdateProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.GenericResponse{Message: "invalid id"})
	}

	var product models.Product
	product.Id = uint(id)

	if err = c.Bind(&product); err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)

	var savedProduct models.Product
	database.DB.Where("id = ?", id).First(&savedProduct)

	return c.JSON(http.StatusOK, savedProduct)
}

func DeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.GenericResponse{Message: "invalid id"})
	}

	database.DB.Model(&models.Product{}).Delete("id = ?", id)
	return c.JSON(http.StatusOK, helpers.GenericResponse{Message: "success"})
}
