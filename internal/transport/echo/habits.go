package echo

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/transport/types"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) HabitsGetAll(c echo.Context) error {
	// If there is a day, then filter by day
	// Check also for date ranges
	// else return the whole thing

	hs, err := h.App.HabitsGetAll(time.Now(), time.Now())
	if err != nil {
		return err
	}

	habits := []*types.HabitTransport{}
	for _, habit := range hs {
		habits = append(habits, types.FromHabitToTransport(habit))
	}
	return c.JSON(http.StatusOK, habits)
}

func (h *Handler) HabitGetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	habit, err := h.App.HabitGetByID(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.FromHabitToTransport(habit))
}

func (h *Handler) HabitCreate(c echo.Context) error {
	day := c.Get("day").(*entity.Day)

	data := &types.HabitTransport{}
	if err := c.Bind(data); err != nil {
		return err
	}

	habit, err := h.App.HabitCreate(day, data.ToHabit())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, types.FromHabitToTransport(habit))
}

func (h *Handler) HabitsGetAllCategories(c echo.Context) error {
	categories, err := h.App.HabitsGetCategories()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, categories)
}
