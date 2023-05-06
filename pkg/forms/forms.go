package forms

import (
	"fmt"
	"net/url"
	"regexp"
)

var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) MinLength(field string, requiredLength int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if len(value) < requiredLength {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", requiredLength))
	}

}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, fmt.Sprint("Pattern doesnt match"))
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
