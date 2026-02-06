package ui

import (
	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func layout(title string, s *State, nodes ...Node) Node {
	if title == "" {
		title = "Curious Ape"
	}
	return HTML5(HTML5Props{
		Title:    "Curious Ape - " + title,
		Language: "en",
		Head: []Node{
			Script(Type("module"), Src("https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js")),
			Link(Rel("stylesheet"), Href("/assets/css/main.css")),
		},
		Body: []Node{
			Class(cLayout),
			Header(
				H1(Text(title)),
			),
			Aside(
				If(s.Authenticated, Nav(
					// TODO: Make pages from nav be partially loaded, and not the full page.
					a("/", "Home "),
					a("/sleep", "Sleep "),
					a("/fitness", "Fitness"),
					a("/deep_work", "Deep-Work "),
					a("/integrations", "Integrations"),
				)),
			),
			Main(
				Group(nodes),
				If(s.Authenticated,
					Button(Text("Logout"), Class("button"), ds.On("click", "@delete('/login')")),
				)),
			Footer(
				P(Text(s.Version)),
			),
		},
	})
}
