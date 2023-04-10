package web

import (
	"html/template"

	"github.com/danielcosme/curious-ape/internal/core/application"
)

type Handler struct {
	App           *application.App
	templateCache map[string]*template.Template
}
