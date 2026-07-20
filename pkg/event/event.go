package event

// type Listener interface {
// 	Listen(data any) error
// }

type Listener func(data any) error
