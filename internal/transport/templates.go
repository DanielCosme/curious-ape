package transport

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/application"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/web"
)

var functions = template.FuncMap{
	"humanDate": humanDate,
	"dateOnly":  dateOnly,
	"lastMonth": lastMonth,
}

type templateData struct {
	CurrentYear     int
	Version         string
	Habit           *core.Habit
	Habits          []*core.Habit
	Days            []dayContainer
	Day             *dayContainer
	Form            any
	Flash           string
	IsAuthenticated bool
	Integrations    []application.IntegrationInfo
}

func (t *Transport) Render(w io.Writer, name string, data any, c echo.Context) error {
	if strings.HasPrefix(name, "p.") {
		ts, ok := t.partialTemplateCache[strings.TrimPrefix(name, "p.")]
		if !ok {
			return fmt.Errorf("the partial template %s does not exist", name)
		}
		return ts.Execute(w, data)
	}

	ts, ok := t.templateCache[name]
	if !ok {
		return fmt.Errorf("the template %s does not exist", name)
	}
	return ts.ExecuteTemplate(w, "base", data)
}

func (t *Transport) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           t.SessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: t.IsAuthenticated(r),
		Version:         t.Version,
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
			// "html/base.gohtml",
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

func humanDate(t time.Time) string {
	return t.Format(core.HumanDate)
}

func dateOnly(t time.Time) string {
	return t.Format(time.DateOnly)
}

func lastMonth(t time.Time) string {
	d := t.AddDate(0, -1, 0)
	d = time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location()).
		AddDate(0, 1, -1)
	return dateOnly(d)
}
