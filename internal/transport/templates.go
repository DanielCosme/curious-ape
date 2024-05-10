package transport

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
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

func (t *Transport) Render(w io.Writer, name string, data any, c echo.Context) error {
	if strings.HasPrefix(name, "-partial-") {
		ts, ok := t.partialTemplateCache[strings.TrimPrefix(name, "-partial-")]
		if !ok {
			return fmt.Errorf("the template %s does not exist", name)
		}
		return ts.Execute(w, data)
	}

	ts, ok := t.templateCache[name]
	if !ok {
		return fmt.Errorf("the template %s does not exist", name)
	}
	return ts.ExecuteTemplate(w, "base", data)
}

func partial(p string) string {
	return "-partial-" + p
}

func humanDate(t time.Time) string {
	return t.Format(entity2.HumanDate)
}

func dateOnly(t time.Time) string {
	return t.Format(time.DateOnly)
}
