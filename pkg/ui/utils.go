package ui

import (
	"fmt"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func GetNextPrevButtons(date core.Date, route string) (prev, next Node) {
	p, n := GetNextPrev(date, route)
	prev = Button(Class(CBtnNav), Text("Previous Month"), ds.On("click", p))
	// TOOD: If next month is in the future don't show button (return nil), or disable it...
	next = Button(Class(CBtnNav), Text("Next Month"), ds.On("click", n))
	return
}

func GetNextPrev(date core.Date, route string) (prev, next string) {
	p, n := GetNextAndPreviousMonth(date)
	prev = fmt.Sprintf("@get('/%s?date=%s')", route, p)
	next = fmt.Sprintf("@get('/%s?date=%s')", route, n)
	return
}

func GetNextAndPreviousMonth(date core.Date) (prev, next string) {
	t := date.FirstDayOfTheMonth().Time()
	previousMonth := t.AddDate(0, -1, 0)
	nextMonth := t.AddDate(0, 1, 0)
	now := time.Now()
	if previousMonth.Month() == now.Month() {
		previousMonth = now
	} else if nextMonth.Month() == now.Month() {
		nextMonth = now
	}
	prev = core.TimeFormatISO8601(previousMonth)
	next = core.TimeFormatISO8601(nextMonth)
	return
}
