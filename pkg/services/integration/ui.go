package integration

import (
	"fmt"
	"strings"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	"danicos.dev/daniel/curious-ape/web/resources"
	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func UI_Integrations(s *ui.UIState, integrations []core.IntegrationInfo) Node {
	return ui.UILayout(s, Div(
		Class(ui.CSurface),
		Map(integrations, func(i core.IntegrationInfo) Node {
			return ui_integration(i)
		}),
	))
}

func ui_integration(i core.IntegrationInfo) Node {
	integrationName := strings.ToLower(i.Name)
	onLoad := fmt.Sprintf("@get('/integrations/%s')", integrationName)
	return Article(
		Class("integration"),
		ID("itg-"+integrationName),
		ds.Init(onLoad),
		H3(
			Img(
				Src(resources.StaticPath(fmt.Sprintf("assets/icons/%s.svg", integrationName))), Alt(i.Name), Class("integration-icon"),
			),
			Text(i.Name),
		),
		P(
			Div(Class(ui.CStatusBadge+" status-"+string(i.Status)), Text(string(i.Status))),
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
				Button(Class(ui.CBtn), Text("Authenticate")),
			),
		),
	)
}
