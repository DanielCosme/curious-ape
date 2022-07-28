package routes

import "github.com/danielcosme/curious-ape/internal/core/application"

type Handler struct {
	App *application.App
}

type envelope map[string]interface{}

func envelopeSuccess() envelope {
	return envelope{"success": "ok"}
}
