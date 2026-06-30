package day

import (
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/dove"
	"danicos.dev/daniel/curious-ape/pkg/ui"
)

func (a *App) HandleDaysMonth(c *dove.Context) (err error) {
	paramDate := getDateParam(c)
	days, err := a.Month(paramDate, core.DESC)
	if err == nil {
		s := ui.StateFromContext(c.Ctx())
		l := ui.Layout("Days", s, UI_days(days))
		return c.RenderOK(l)
	}
	return err
}

func getDateParam(c *dove.Context) core.Date {
	c.ParseForm()
	if c.Req.Form.Get("date") == "" {
		return core.NewDate(time.Now())
	} else {
		date, err := core.NewDateFromISO8601(c.Req.Form.Get("date"))
		if err == nil {
			return date
		}
		c.Log.Fatal("cannot parse date", "err", err)
		panic(err)
	}
}
