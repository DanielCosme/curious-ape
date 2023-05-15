package entity

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
	GuestRole Role = "guest"
)

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
	Email    string `db:"email"`
	Role     Role   `db:"role"`
}

type UserFilter struct {
	ID       int
	Role     Role
	Password string
}
