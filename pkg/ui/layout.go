package ui

import (
	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/web/resources"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func UILayout(s UIState, nodes ...Node) Node {
	if s.Title == "" {
		s.Title = "Curious Ape"
	}
	return HTML5(HTML5Props{
		Title:    "Curious Ape - " + s.Title,
		Language: "en",
		Head: []Node{
			Script(Defer(), Type("module"), Src(resources.StaticPath("datastar/datastar.js"))),
			Link(Rel("stylesheet"), Href(resources.StaticPath("main.css")), Type("text/css")),

			// <link rel="icon" type="image/x-icon" href={ resources.StaticPath("assets/favicon.ico") }/>

			// Link(Rel("preload"), Href("/assets/fonts/Exo2-Regular.woff2"), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			// Link(Rel("preload"), Href("/assets/fonts/Exo2-SemiBold.woff2"), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			// Link(Rel("preload"), Href("/assets/fonts/FiraCode-Regular.woff2"), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
			// Link(Rel("preload"), Href("/assets/fonts/FiraCode-Bold.woff2"), As("font"), Type("font/woff2"), CrossOrigin("anonymous")),
		},
		Body: []Node{
			Class(cLayout),
			Header(
				H1(Text(s.Title)),
			),
			Group(nodes),
		},
	})
}

func Layout(title string, s *State, nodes ...Node) Node {
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
			Group(nodes),
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
