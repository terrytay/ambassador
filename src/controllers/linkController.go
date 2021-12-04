package controllers

import (
	"net/http"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/labstack/echo/v4"
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/helpers"
	"github.com/terrytay/ambassador/src/middlewares"
	"github.com/terrytay/ambassador/src/models"
)

func Link(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.GenericResponse{Message: "invalid id"})
	}

	var links []models.Link

	database.DB.Where("user_id = ?", id).Find(&links)

	for i, link := range links {
		var orders []models.Order

		database.DB.Where("code = ? and complete = true", link.Code).Find(&orders)

		links[i].Orders = orders
	}

	return c.JSON(http.StatusOK, links)
}

type CreateLinkRequest struct {
	Products []int
}

func CreateLink(c echo.Context) error {
	var request CreateLinkRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err,
		})
	}

	id, _ := middlewares.GetUserId(c)

	link := models.Link{
		UserId: id,
		Code:   faker.Username(),
	}

	for _, productId := range request.Products {
		product := models.Product{}
		product.Id = uint(productId)
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)
	database.DB.Preload("Products").Find(&link)
	database.DB.Preload("User").Find(&link)

	return c.JSON(http.StatusOK, link)
}
