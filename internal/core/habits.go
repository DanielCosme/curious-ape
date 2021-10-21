package core

type Habit struct {
	ID     int    `json:"id"`
	State  string `json:"state"`
	Date   string `json:"date"`
	Origin string `json:"origin"`
	Type   string `json:"type"`
}

type HabitModel interface {
	GetAll() ([]Habit, error)
	Insert(h *Habit) error
	Get(id int) (*Habit, error)
	UpdateOrCreate(h *Habit) error
	Update(h *Habit) error
	UpdateByDate(h *Habit) error
	Delete(id int) error
}
