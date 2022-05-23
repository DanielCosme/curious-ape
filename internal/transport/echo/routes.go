package echo

import (
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/transport/echo/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"net/http"
)

func Routes(a *application.App) http.Handler {
	h := Handler{App: a}
	e := echo.New()

	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())

	e.GET("/ping", h.Ping)

	days := e.Group("/days")
	{
		days.GET("", h.DaysGetAll)
		daysByDate := days.Group("/:date", middleware.SetDay(a))
		{ // /days/:date/habits
			daysByDate.POST("/habits", h.HabitCreate)
		}
	}

	habits := e.Group("/habits")
	{
		habits.GET("", h.HabitsGetAll)
		habits.POST("", h.HabitCreate, middleware.SetDay(a))
		habits.GET("/categories", h.HabitsGetAllCategories)
		habitsByID := habits.Group("/:id", middleware.SetHabit(a))
		{
			habitsByID.GET("", h.HabitGetByID)
		}
	}

	return e
}
