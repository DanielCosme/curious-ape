package ui

import (
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
		H2(Text("2026")),
		deadline(),
		deadline(),
	))
}

func deadline() Node {
	return Section(
		H4(Text("Laura's Birthday")),
		P(Text("7th July")),
		P(Text("Days left: 50")),
		P(Text("Percentage: 30%")),
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
			Input(
				Type("test"), Name("title"), ds.Bind("name"),
			),
			Input(
				Type("date"),
				Name("end_date"),
				ds.Bind("end_date"),
			),
			Input(
				Type("checkbox"),
				Name("recurrent"),
				ds.Bind("recurrent"),
			),
			Button(
				Text("Create Deadline"),
			),
		),
	))
}
