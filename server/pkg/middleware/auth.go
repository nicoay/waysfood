package middleware

import (
	dto "mytask/dto/result"
	jwtToken "mytask/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Declare Result struct here ...
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Create Auth here
func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization") //Take value from header authorization from request HTTP used func c.Req() from echo.Context

		if token == "" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusBadRequest, Message: "unauthorized"})
		}

		token = strings.Split(token, " ")[1]
		claims, err := jwtToken.DecodeToken(token)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, Result{Code: http.StatusUnauthorized, Message: "unathorized"})
		}

		c.Set("userLogin", claims)
		return next(c)
	}
}
