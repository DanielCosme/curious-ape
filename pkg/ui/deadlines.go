package ui

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Deadlines(s *State) Node {
	return layout("Deadlines", s, Div(
		Form(
			// This redirects to the form to create a new deadline.
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
