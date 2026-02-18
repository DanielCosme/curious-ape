package ui

import (
	"fmt"
	"net/url"
	"strings"

	"git.danicos.dev/daniel/curious-ape/pkg/application"
	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func Integrations(s *State) Node {
	return layout("Integrations", s, Div(
		Map(s.Integrations, func(i application.IntegrationInfo) Node {
			return Integration(i)
		}),
	))
}

func Integration(i application.IntegrationInfo) Node {
	integrationName := strings.ToLower(i.Name)
	q := url.Values{}
	q.Add("name", integrationName)
	onLoad := fmt.Sprintf("@get('/integration?%s')", q.Encode())
	return Article(
		ID("itg-"+integrationName),
		ds.Init(onLoad),
		H3(Text(i.Name)),
		P(Text(fmt.Sprintf("Status: %s", i.Status))),
		If(len(i.Info) > 0,
			Ul(
				Map(i.Info, func(info string) Node {
					return Li(Text(info))
				}),
			),
		),
		If(i.Status == application.IntegrationStatusDisconnected && i.AuthURL != "",
			A(
				Href(i.AuthURL),
				Target("_blank"),
				Button(Text("Authenticate")),
			),
		),
	)
}
