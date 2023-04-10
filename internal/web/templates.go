package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/entity"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
}

type templateData struct {
	CurrentYear int
	Habit       *entity.Habit
	Habits      []*entity.Habit
}

func (h *Handler) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (h *Handler) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := h.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func humanDate(t time.Time) string {
	return t.Format(entity.HumanDate)
}
