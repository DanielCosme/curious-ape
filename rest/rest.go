package rest

type Envelope map[string]interface{}

func Payload(title string, payload interface{}) *Envelope {
	return &Envelope{title: payload}
}

func Message(payload interface{}) *Envelope {
	return Payload("message", payload)
}
