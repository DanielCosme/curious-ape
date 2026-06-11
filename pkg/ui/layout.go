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
			// Preload critical Exo 2 (sans) + Fira Code (mono) weights.
			// Mono preloads help the prominent data/scores/grids appear fast.
			Link(Rel("preload"), Href("/assets/fonts/Exo2-Regular.woff2"), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			Link(Rel("preload"), Href("/assets/fonts/Exo2-SemiBold.woff2"), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			Link(Rel("preload"), Href("/assets/fonts/FiraCode-Regular.woff2"), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			Link(Rel("preload"), Href("/assets/fonts/FiraCode-Bold.woff2"), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
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
					navItem(lucide.House(), "/", "Home ", s.CurrentPath),
					navItem(lucide.SquareCheckBig(), "/habits", "Habits", s.CurrentPath),
					navItem(lucide.Hourglass(), "/deadlines", "Deadlines ", s.CurrentPath),
					navItem(lucide.Bed(), "/sleep", "Sleep ", s.CurrentPath),
					navItem(lucide.Dumbbell(), "/fitness", "Fitness", s.CurrentPath),
					navItem(lucide.MonitorCog(), "/deep_work", "Deep-Work ", s.CurrentPath),
					navItem(lucide.Workflow(), "/integrations", "Integrations", s.CurrentPath),
				)),
			),
			Main(
				Group(nodes),
				If(s.Authenticated,
					Button(Text("Logout"), Class("btn btn-secondary"), ds.On("click", "@delete('/login')")),
				)),
			Footer(
				P(Class(cVersion), Text(s.Version)),
			),
		},
	})
}

func navItem(icon Node, path, name, current string) Node {
	cls := cNavLink
	if path == current {
		cls = cNavLinkActive
	}
	return A(
		Class(cls),
		Href(path),
		icon, Text(name),
	)
}
