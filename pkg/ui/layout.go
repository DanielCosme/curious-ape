package ui

import (
	"danicos.dev/daniel/curious-ape/pkg/config"
	lucide "github.com/eduardolat/gomponents-lucide"
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
			Script(Type("module"), Src(config.DATASTAR)),
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
					navItem(lucide.House(), "/", "Home "),
					navItem(lucide.SquareCheckBig(), "/habits", "Habits"),
					navItem(lucide.Hourglass(), "/deadlines", "Deadlines "),
					navItem(lucide.Bed(), "/sleep", "Sleep "),
					navItem(lucide.Dumbbell(), "/fitness", "Fitness"),
					navItem(lucide.MonitorCog(), "/deep_work", "Deep-Work "),
					navItem(lucide.Workflow(), "/integrations", "Integrations"),
				)),
			),
			Main(
				Group(nodes),
				If(s.Authenticated,
					Button(Text("Logout"), Class("btn btn-secondary"), ds.On("click", "@delete('/login')")),
				)),
			Footer(
				P(Text(s.Version)),
			),
		},
	})
}

func navItem(icon Node, path, name string) Node {
	return A(
		Class(cNavLink),
		icon, Href(path), Text(name),
	)
}
