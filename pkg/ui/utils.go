package ui

import (
	"fmt"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func GetNextPrevButtons(day core.Day, route string) (prev, next Node) {
	p, n := GetNextPrev(day, route)
	prev = Button(Class(cBtnNav), Text("Previous Month"), ds.On("click", p))
	next = Button(Class(cBtnNav), Text("Next Month"), ds.On("click", n))
	return
}

func GetNextPrev(day core.Day, route string) (prev, next string) {
	p, n := GetNextAndPreviousMonth(day)
	prev = fmt.Sprintf("@get('/%s?date=%s')", route, p)
	next = fmt.Sprintf("@get('/%s?date=%s')", route, n)
	return
}

func GetNextAndPreviousMonth(day core.Day) (prev, next string) {
	t := day.Date.FirstDayOfTheMonth().Time()
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
