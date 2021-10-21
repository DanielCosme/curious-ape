package errors

var (
	ErrRecordNotFound = New("db record not found")
	ErrEditConflict   = New("edit conflict")
	ErrDuplicateEmail = New("duplicate email")
)
