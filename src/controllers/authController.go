package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/terrytay/ambassador/src/database"
	"github.com/terrytay/ambassador/src/middlewares"
	"github.com/terrytay/ambassador/src/models"
)

func Register(c echo.Context) error {
	var data map[string]string

	if err := c.Bind(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "passwords do not match",
		})
	}

	var emailAlreadyExists string
	database.DB.Model(&models.User{}).Where("email = ?", data["email"]).Pluck("email", &emailAlreadyExists)

	if emailAlreadyExists != "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "email in use",
		})
	}

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: false,
	}
	if err := user.SetPassword(data["password"]); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "please try again later",
		})
	}

	database.DB.Create(&user)

	return c.JSON(http.StatusOK, user)
}

func Login(c echo.Context) error {
	var data map[string]string

	if err := c.Bind(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid credentials",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid credentials",
		})

	}

	expiresAt := time.Now().Add(time.Hour * 24)

	payload := jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.Id)),
		ExpiresAt: expiresAt.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid credentials",
		})
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  expiresAt,
		HttpOnly: true,
	}

	c.SetCookie(&cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}

func User(c echo.Context) error {
	id, _ := middlewares.GetUserId(c)

	var user models.User
	database.DB.Where("id = ?", id).First(&user)

	return c.JSON(http.StatusOK, user)
}

func Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "jwt", Value: ""})
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}

func UpdateInfo(c echo.Context) error {
	var data map[string]string

	if err := c.Bind(&data); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(c)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}
	var emailAlreadyExists string
	database.DB.Model(&models.User{}).Where("email = ? AND id <> ?", data["email"], id).Pluck("email", &emailAlreadyExists)

	if emailAlreadyExists != "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "email in use",
		})
	}

	database.DB.Model(models.User{}).Where("id = ?", id).Updates(&user)

	return c.JSON(http.StatusOK, user)
}

func UpdatePassword(c echo.Context) error {
	var data map[string]string

	if err := c.Bind(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "password do not match",
		})
	}

	id, _ := middlewares.GetUserId(c)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	if err := user.SetPassword(data["password"]); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "please try again later",
		})
	}

	database.DB.Updates(&user)
	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}
