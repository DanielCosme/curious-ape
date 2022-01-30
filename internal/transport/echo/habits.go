package echo

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) HabitsGetAll(c echo.Context) error {
	habits, err := h.App.Habits.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, habits)
}

func (h *Handler) HabitsGetByID(c echo.Context) error {
	id := c.Param("id")
	habits, err := h.App.Habits.GetByID(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, habits)
}

func (h *Handler) HabitCreate(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
