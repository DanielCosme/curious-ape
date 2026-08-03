package deadline

import (
	"fmt"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func UI_Deadlines(s *ui.UIState, deadlineState ui.DeadlineState) Node {
	return ui.UILayout(s, Div(
		Class(ui.CSurface),
		A(
			Href("/deadlines/new"),
			Button(Class(ui.CBtn), Text("New deadline")),
		),
		Map(deadlineState.DS, func(d core.Deadline) Node {
			return deadline(d)
		}),
	))
}

func deadline(d core.Deadline) Node {
	if d.EndDate.Time.IsZero() {
		return nil
	}
	return Div(
		Class(ui.CLogEntry+" deadline-item"),
		H4(Text(d.Title)),
		P(Text(d.EndDate.Time.Format("02 Jan 2006"))),
		P(Text(fmt.Sprintf("Days left: %d", d.DaysLeft))),
	)
}

func DeadlineForm(s *ui.UIState, deadlineState ui.DeadlineState) Node {
	post := "@post('/deadlines/new', {contentType: 'form'})"

	var err string
	if deadlineState.Err != nil {
		err = deadlineState.Err.Error()
	}
	return ui.UILayout(s, Div(
		Class(ui.CSurface),
		If(deadlineState.Err != nil,
			P(Class(ui.CError), Text("ERROR: "+err)),
		),
		Form(
			ds.On("submit", post),
			Label(
				Text("Title"),
				For("title"),
				Input(Type("text"), Name("title"), ds.Bind("name")),
			),
			Label(
				For("end_date"),
				Text("End Date"),
				Input(
					Type("date"),
					Name("end_date"),
					ds.Bind("end_date"),
				),
			),
			Label(
				For("recurrent"),
				Text("Recurrent"),
				Input(
					Type("checkbox"),
					Name("recurrent"),
					ds.Bind("recurrent"),
				),
			),
			Button(
				Class(ui.CBtn),
				Text("Create Deadline"),
			),
		),
	))
}
