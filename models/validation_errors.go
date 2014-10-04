package models

import "fmt"

type ValidationErrors struct {
	Errors map[string][]string
}

func doThing() {
	things := ValidationErrors{}
	fmt.Println(things)
}

func NewValidationErrors() ValidationErrors {
	return ValidationErrors{
		make(map[string][]string),
	}
}

func (r ValidationErrors) Error() string {
	return "error"
}

func (r *ValidationErrors) Add(name, message string) {
	if r.Errors[name] == nil {
		r.Errors[name] = []string{}
	}

	r.Errors[name] = append(r.Errors[name], message)
}

func (r ValidationErrors) Any() bool {
	return len(r.Errors) > 0
}
