package user

import (
	"danicos.dev/daniel/curious-ape/pkg/ui"
	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func UI_Login(s ui.UIState) Node {
	form := Form(
		Class("login-form"),
		FieldSet(
			Legend(Text("Login")),
			Div(
				Label(For("username"), Text("Username")),
				Input(Type("text"), ID("username"), Name("username"), Placeholder(""), Required()),
			),
			Div(
				Label(For("password"), Text("Password")),
				Input(Type("text"), ID("password"), Name("password"), Placeholder(""), Required()),
			),
			Button(Class(ui.CBtn),
				Text("Login"),
				Type("submit"),
				ds.On("click", "@post('/login', {contentType: 'form'})"),
			),
		),
	)
	s.Title = "Login"
	return ui.UILayout(s, form)
}
