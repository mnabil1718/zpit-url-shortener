package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func ErrorHandler(c *echo.Context, err error) {

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	htmxReq := c.Request().Header.Get("HX-Request")
	if htmxReq == "true" {
		c.Render(code, "error-message", err.Error())
	} else {
		c.String(code, err.Error())
	}

}
