package transport

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Day struct {
	ID   int32  `json:"id"`
	Date string `json:"date"`
}

func (t *Transport) home(c echo.Context) error {
	days, err := t.App.DaysMonth(core.NewDateToday())
	if err != nil {
		return errServer(err)
	}

	return c.JSON(http.StatusOK, days)
}

// func dayToTranspor(days []*core.Day)
