package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/models"
)

func Orders(c echo.Context) error {
	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)

	for i, order := range orders {
		orders[i].Name = order.FullName()
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(http.StatusOK, orders)
}

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrder(c echo.Context) error {
	var request CreateOrderRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err,
		})
	}

	link := models.Link{}

	database.DB.Preload("User").First(&link, models.Link{
		Code: request.Code,
	})

	if link.Id == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid link",
		})
	}

	order := models.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: link.User.Email,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err,
		})
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, requestProduct := range request.Products {
		product := models.Product{}
		product.Id = uint(requestProduct["product_id"])

		tx.First(&product)

		total := product.Price * float64(requestProduct["quantity"])

		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err,
			})
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(product.Title),
			Description: stripe.String(product.Description),
			Images:      []*string{stripe.String(product.Image)},
			Amount:      stripe.Int64(int64(product.Price) * 100),
			Currency:    stripe.String("sgd"),
			Quantity:    stripe.Int64(int64(requestProduct["quantity"])),
		})
	}

	stripe.Key = "sk_test_51K3kZ4CJjXVU2nWfPBXysSYnS4JZxTvd8NmxQcdpR1YM5OGsnyjG0tHXwUY40GamddqcZzICpWgoaHgYsgW9EK6F00USvXuv68"

	params := stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String("http://localhost:5000/success?source={CHECKOUT_SESSION_ID}"),
		CancelURL:          stripe.String("http://localhost:5000/error"),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
	}

	source, err := session.New(&params)
	if err != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err,
		})
	}

	order.TransactionId = source.ID

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err,
		})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, source)
}

func CompleteOrder(c echo.Context) error {
	var data map[string]string

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err,
		})
	}

	order := models.Order{}

	database.DB.Preload("OrderItems").First(&order, models.Order{
		TransactionId: data["source"],
	})

	if order.Id == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "order not found",
		})
	}

	order.Complete = true

	database.DB.Save(&order)

	go func(order models.Order) {
		ambassadorRevenue := 0.0
		adminRevenue := 0.0

		for _, item := range order.OrderItems {
			ambassadorRevenue += item.AmbassadorRevenue
			adminRevenue += item.AdminRevenue
		}

		user := models.User{}
		user.Id = order.UserId

		database.DB.Find(&user)

		database.Cache.ZIncrBy(context.Background(), "rankings", ambassadorRevenue, user.Name())

		ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link %s", ambassadorRevenue, order.Code))
		smtp.SendMail("mailhog:1025", nil, "no-reply@ambassadors.com", []string{order.AmbassadorEmail}, ambassadorMessage)

		adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed", order.Id, adminRevenue))
		smtp.SendMail("mailhog:1025", nil, "no-reply@ambassadors.com", []string{"txrrythk@gmail.com"}, adminMessage)
	}(order)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}
