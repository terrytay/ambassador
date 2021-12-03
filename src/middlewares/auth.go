package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const SecretKey = "secret"

type ClaimsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("jwt")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "unauthenticated",
			})
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "unauthenticated",
			})
		}

		payload := token.Claims.(*ClaimsWithScope)
		isAmbassadorPath := strings.Contains(c.Path(), "/api/ambassador")

		if (payload.Scope == "admin" && isAmbassadorPath) || (payload.Scope == "ambassador" && !isAmbassadorPath) {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"message": "unauthorized",
			})
		}

		return next(c)
	}

}

func GenerateJWT(id uint, expiresAt time.Time, scope string) (string, error) {

	payload := ClaimsWithScope{}
	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = expiresAt.Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(SecretKey))

}

func GetUserId(c echo.Context) (uint, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return 0, c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthenticated",
		})
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		return 0, c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthenticated",
		})
	}

	payload := token.Claims.(*ClaimsWithScope)
	id, _ := strconv.Atoi(payload.Subject)

	return uint(id), nil
}
