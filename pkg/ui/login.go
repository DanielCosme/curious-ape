package ui

import (
	. "maragu.dev/gomponents"
	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func Login(s *State) Node {
	h1 := H1(Text("Login"))
	form := Form(
		FieldSet(
			Legend(Text("Login")),
			Div(
				Label(For("username"), Text("Username")),
				Input(Type("text"), ID("username"), Name("username"), Placeholder(""), Required(), blockDisplay()),
			),
			Div(
				Label(For("password"), Text("Password")),
				Input(Type("text"), ID("password"), Name("password"), Placeholder(""), Required(), blockDisplay()),
			),
			Button(Text("Login"),
				Type("submit"),
				ds.On("click", "@post('/login', {contentType: 'form'})"),
			),
		),
	)
	return layout("login", s, h1, form)
}
