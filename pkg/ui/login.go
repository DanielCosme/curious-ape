package ui

import (
	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func Login(s *State) Node {
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
			Button(Class(cBtn),
				Text("Login"),
				Type("submit"),
				ds.On("click", "@post('/login', {contentType: 'form'})"),
			),
		),
	)
	return layout("Login", s, form)
}
