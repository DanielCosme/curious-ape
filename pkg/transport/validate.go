package transport

import "fmt"

type Validation interface {
	Validate(v *Validator)
}

type ValidationError struct {
	v *Validator
}

func (e *ValidationError) Error() string {
	return e.v.message
}

type Validator struct {
	message string
	details map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		details: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.details) == 0
}

func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.Add(field, message)
	}
}

func (v *Validator) Add(field, message string) {
	v.details[field] = message
}

func (v *Validator) Message(msg string) {
	if !v.Valid() && v.message == "" {
		v.message = msg
	}
}

func invalid(s string) string {
	if s == "" {
		return "empty"
	}
	return fmt.Sprintf("invalid value: %v", s)
}
