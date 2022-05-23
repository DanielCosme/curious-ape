package middleware

import (
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/labstack/echo/v4"
	"strconv"
)

func SetDay(a *application.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.Request().Header.Get("X-APE-DATE")
			if key == "" {
				key = c.Param("date")
			}

			date, err := entity.ParseDate(key)
			if err != nil {
				return err
			}

			day, err := a.DayGetByDate(date)
			if err != nil {
				return err
			}

			c.Set("day", day)
			return next(c)
		}
	}
}

func SetHabit(a *application.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				return err
			}

			habit, err := a.HabitGetByID(id)
			if err != nil {
				return err
			}

			c.Set("habit", habit)
			return next(c)
		}
	}
}
