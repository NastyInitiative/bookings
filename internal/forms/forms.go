package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form - creates a custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New - initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has - checks if form field in post and not empty
func (form *Form) Has(field string) bool {
	reqField := form.Get(field)
	if reqField == "" {
		return false
	}
	return true
	// return reqField != ""
}

// Valid - return true if there are no errors, otherwise false
func (form *Form) Valid() bool {
	return len(form.Errors) == 0
}

// Required - checks if fields are empty. If they are, an error message will shown.
func (form *Form) Required(fields ...string) {
	for _, field := range fields {
		value := form.Get(field)
		if strings.TrimSpace(value) == "" {
			form.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MinLength - checks for string minimum length
func (form *Form) MinLength(field string, length int) bool {
	reqField := form.Get(field)
	if len(reqField) < length {
		form.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail - Checks for invalid email address
func (form *Form) IsEmail(field string) {
	if !govalidator.IsEmail(form.Get(field)) {
		form.Errors.Add(field, "Invalid email address")
	}
}
