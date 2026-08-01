package ui

import (
	"fmt"

	"danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func Deadlines(s *State) Node {
	return Layout("Deadlines", s, Div(
		Class(CSurface),
		A(
			Href("/deadline"),
			Button(Class(CBtn), Text("New deadline")),
		),
		Map(s.Deadlines.DS, func(d core.Deadline) Node {
			return deadline(d)
		}),
	))
}

func deadline(d core.Deadline) Node {
	if d.EndDate.Time().IsZero() {
		return nil
	}
	return Div(
		Class(CLogEntry+" deadline-item"),
		H4(Text(d.Title)),
		P(Text(d.EndDate.Time.Format("02 Jan 2006"))),
		P(Text(fmt.Sprintf("Days left: %d", d.DaysLeft))),
	)
}

func DeadlineForm(s *State) Node {
	post := "@post('/deadline', {contentType: 'form'})"

	var err string
	if s.Deadlines.Err != nil {
		err = s.Deadlines.Err.Error()
	}
	return Layout("New Deadline", s, Div(
		Class(CSurface),
		If(s.Deadlines.Err != nil,
			P(Class(CError), Text("ERROR: "+err)),
		),
		Form(
			ds.On("submit", post),
			Label(
				Text("Title"),
				For("title"),
				Input(Type("text"), Name("title"), ds.Bind("name")),
			),
			Label(
				For("end_date"),
				Text("End Date"),
				Input(
					Type("date"),
					Name("end_date"),
					ds.Bind("end_date"),
				),
			),
			Label(
				For("recurrent"),
				Text("Recurrent"),
				Input(
					Type("checkbox"),
					Name("recurrent"),
					ds.Bind("recurrent"),
				),
			),
			Button(
				Class(CBtn),
				Text("Create Deadline"),
			),
		),
	))
}
