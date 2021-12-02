package middlewares

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "unauthenticated",
			})
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "unauthenticated",
			})
		}

		return next(c)
	}

}

func GetUserId(c echo.Context) (uint, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return 0, c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthenticated",
		})
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return 0, c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthenticated",
		})
	}

	payload := token.Claims.(*jwt.StandardClaims)
	id, _ := strconv.Atoi(payload.Subject)

	return uint(id), nil
}
