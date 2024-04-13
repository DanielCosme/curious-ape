package transport

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	entity2 "github.com/danielcosme/curious-ape/internal/entity"

	"github.com/danielcosme/curious-ape/web"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
	"dateOnly":  dateOnly,
}

type templateData struct {
	CurrentYear     int
	Version         string
	Habit           *entity2.Habit
	Habits          []*entity2.Habit
	Days            []dayContainer
	Day             *dayContainer
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func (h *Transport) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           h.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: h.IsAuthenticated(r),
		Version:         h.Version,
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(web.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.gohtml",
			"html/partials/*gohtml",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(web.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func newTemplatePartialCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(web.Files, "html/partials/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.gohtml",
			"html/partials/*gohtml",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(web.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

func (h *Transport) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := h.templateCache[page]
	if !ok {
		h.serverError(w, fmt.Errorf("the template %s does not exist", page))
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.serverError(w, err)
		return
	}

	if _, err := buf.WriteTo(w); err != nil {
		panic(err)
	}
	w.WriteHeader(status)
}

func (h *Transport) renderPartial(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := h.partialTemplateCache[page]
	if !ok {
		http.Error(w, fmt.Sprintf("the template %s does not exist", page), http.StatusInternalServerError)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := buf.WriteTo(w); err != nil {
		panic(err)
	}
	w.WriteHeader(status)
}

func humanDate(t time.Time) string {
	return t.Format(entity2.HumanDate)
}

func dateOnly(t time.Time) string {
	return t.Format(time.DateOnly)
}
