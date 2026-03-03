package main

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v5"
	qrlib "github.com/mnabil1718/zp.it/internal/qr"
	"github.com/mnabil1718/zp.it/internal/shortener"
)

func Health(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

func Index(c *echo.Context) error {
	return c.Render(200, "index", nil)
}

type Result struct {
	Short    string
	Original string
	QRCode   string // base64
}

func Generate(c *echo.Context) error {
	url := c.FormValue("url")
	qr := c.FormValue("qr") == "on"

	code, err := shortener.Shorten(6)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed shorten URL")
	}

	s := "http://localhost:8080/" + code

	data := Result{
		Short:    s,
		Original: url,
	}

	if qr {
		png, err := qrlib.GenerateQR(s)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Cannot process QR generation")
		}

		data.QRCode = base64.StdEncoding.EncodeToString(png)
	}

	return c.Render(200, "result", data)
}
