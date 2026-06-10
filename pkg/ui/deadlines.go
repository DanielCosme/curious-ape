package ui

import (
	"fmt"

	"danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func Deadlines(s *State) Node {
	return layout("Deadlines", s, Div(
		A(
			Href("/deadline"),
			Button(Text("New deadline")),
		),
		Map(s.Deadlines.DS, func(d core.Deadline) Node {
			return deadline(d)
		}),
	))
}

func deadline(d core.Deadline) Node {
	return Section(
		H4(Text(d.Title)),
		P(Text(d.EndDate.Time().Format("02 Jan 2006"))),
		P(Text(fmt.Sprintf("Days left: %d", d.DaysLeft))),
	)
}

func DeadlineForm(s *State) Node {
	post := "@post('/deadline', {contentType: 'form'})"

	var err string
	if s.Deadlines.Err != nil {
		err = s.Deadlines.Err.Error()
	}
	return layout("New Deadline", s, Div(
		If(s.Deadlines.Err != nil,
			P(Text("ERROR: "+err)),
		),
		Form(
			ds.On("submit", post),
			Label(
				Style("display: block"),
				Text("Title"),
				For("title"),
				Input(Type("test"), Name("title"), ds.Bind("name")),
			),
			Label(
				Style("display: block"),
				For("end_date"),
				Text("End Date"),
				Input(
					Type("date"),
					Name("end_date"),
					ds.Bind("end_date"),
				),
			),
			Label(
				Style("display: block"),
				For("recurrent"),
				Text("Recurrent"),
				Input(
					Type("checkbox"),
					Name("recurrent"),
					ds.Bind("recurrent"),
				),
			),
			Button(
				Text("Create Deadline"),
			),
		),
	))
}
