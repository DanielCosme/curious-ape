package habit

type State string

const (
	Done    State = "done"
	NoInfo  State = "no-info"
	NotDone State = "not-done"
)

type Code string

const (
	Food       Code = "food"
	DeepWork   Code = "deep-work"
	Fitness    Code = "fitness"
	WakeUp     Code = "wake-up"
	CodeCustom Code = "custom"
)
