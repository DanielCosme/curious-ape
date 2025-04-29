package transport

import (
	"net/http"
)

type Route struct {
	// TODO: Figure out how to handle middlewares.
	pattern string
	handler http.HandlerFunc
	body    any      // What does this route return as a body?
	params  []string // what params does it require and how do we validate this?
}

func Routes(t *Transport) []Route {
	routes := []Route{
		{
			pattern: "GET /api/info",
			handler: t.HandlerInfo,
		},
		{
			pattern: "GET /api/health",
			handler: func(w http.ResponseWriter, r *http.Request) {},
		},
		{
			pattern: "GET /api/days",
			handler: t.HandlerDaysMonth,
		},
		{
			pattern: "POST /api/habits",
			handler: t.HandlerHabitsUpsert,
		},
		{
			pattern: "GET /api/habits/types",
			handler: t.HandlerHabitTypes,
		},
	}
	return routes
}

type Resource struct {
	// TODO: Implement resource abstraction
	// 		- A Resource is a collection of Routes that represent an entity.
	// 			e.g: Habit, Sleep Record, etc...
	// 			- Methods
	// 			- Not Allowed methods
	// 			- Middlewares
	//
	// This abstraction will generate a list of routes that get passed into the "Mux"
	// Also it may be used to validate parameters and generate error messages.

	// TODO Generate information for the options request.
	// TODO Generate 404 405 response for this resource.
}

func (r *Resource) GenerateRoutes() []Route {
	return nil
}
