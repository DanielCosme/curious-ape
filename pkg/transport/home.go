package transport

import (
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type DaysPayload struct {
	Month string       `json:"month"`
	Days  []DaySummary `json:"days"`
}

type DaySummary struct {
	ID     int32        `json:"id"`
	Key    string       `json:"key"`
	Date   time.Time    `json:"date"`
	WakeUp HabitSummary `json:"wake_up"`
	Sleep  HabitSummary `json:"sleep"`
	Work   HabitSummary `json:"work"`
	Eat    HabitSummary `json:"eat"`
}

type HabitSummary struct {
}

func (t *Transport) home(c echo.Context) error {
	days, err := t.App.DaysMonth(core.NewDateToday())
	if err != nil {
		return errServer(err)
	}
	return c.JSON(http.StatusOK, days)
}

// func dayToTranspor(days []*core.DaySummary)
