package ui

import (
	"fmt"
	"net/url"
	"strings"

	"danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func Integrations(s *State) Node {
	return Layout("Integrations", s, Div(
		Class(cSurface),
		Map(s.Integrations, func(i core.IntegrationInfo) Node {
			return Integration(i)
		}),
	))
}

func Integration(i core.IntegrationInfo) Node {
	integrationName := strings.ToLower(i.Name)
	q := url.Values{}
	q.Add("name", integrationName)
	onLoad := fmt.Sprintf("@get('/integration?%s')", q.Encode())
	return Article(
		Class("integration"),
		ID("itg-"+integrationName),
		ds.Init(onLoad),
		H3(
			Img(Src("/assets/icons/"+integrationName+".svg"), Alt(i.Name), Class("integration-icon")),
			Text(i.Name),
		),
		P(
			Div(Class(cStatusBadge+" status-"+string(i.Status)), Text(string(i.Status))),
		),
		If(len(i.Info) > 0,
			Ul(
				Map(i.Info, func(info string) Node {
					return Li(Text(info))
				}),
			),
		),
		If(i.Status == core.IntegrationStatusDisconnected && i.AuthURL != "",
			A(
				Href(i.AuthURL),
				Target("_blank"),
				Button(Class(CBtn), Text("Authenticate")),
			),
		),
	)
}
