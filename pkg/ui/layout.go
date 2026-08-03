package ui

import (
	"danicos.dev/daniel/curious-ape/web/resources"
	lucide "github.com/eduardolat/gomponents-lucide"
	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Layout(s *State, nodes ...Node) Node {
	if s.Title == "" {
		s.Title = "Curious Ape"
	}
	return HTML5(HTML5Props{
		Title:    "Curious Ape - " + s.Title,
		Language: "en",
		Head: []Node{
			Script(Defer(), Type("module"), Src(resources.StaticPath("datastar/datastar.js"))),
			Link(Rel("stylesheet"), Href(resources.StaticPath("main.css")), Type("text/css")),

			Link(Rel("icon"), Type("image/x-icon"), Href(resources.StaticPath("assets/icons/favicon.ico"))),

			Link(Rel("preload"), Href(resources.StaticPath("assets/fonts/Exo2-Regular.woff2")), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			Link(Rel("preload"), Href(resources.StaticPath("assets/fonts/Exo2-SemiBold.woff2")), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			Link(Rel("preload"), Href(resources.StaticPath("assets/fonts/FiraCode-Regular.woff2")), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			Link(Rel("preload"), Href(resources.StaticPath("assets/fonts/FiraCode-Bold.woff2")), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
		},
		Body: []Node{
			Class(CLayout),
			Header(
				H1(Text(s.Title)),
			),
			If(s.IsAuthenticated,
				Aside(
					Nav(
						navItem(lucide.House(), "/", "Home ", s.CurrentPath),
						navItem(lucide.SquareCheckBig(), "/habits", "Habits", s.CurrentPath),
						navItem(lucide.Hourglass(), "/deadlines", "Deadlines ", s.CurrentPath),
						navItem(lucide.Bed(), "/sleep", "Sleep ", s.CurrentPath),
						navItem(lucide.Dumbbell(), "/fitness", "Fitness", s.CurrentPath),
						navItem(lucide.MonitorCog(), "/worklog", "Deep-Work ", s.CurrentPath),
						navItem(lucide.Workflow(), "/integrations", "Integrations", s.CurrentPath),
					),
				),
			),
			Main(
				Group(nodes),
				If(s.IsAuthenticated,
					Button(Text("Logout"), Class("btn btn-secondary"), ds.On("click", "@delete('/login')")),
				),
			),
			Footer(
				P(Class(CVersion), Text(s.Version)),
			),
		},
	})
}

func navItem(icon Node, path, name, current string) Node {
	cls := CNavLink
	if path == current {
		cls = CNavLinkActive
	}
	return A(
		Class(cls),
		Href(path),
		icon, Text(name),
	)
}
