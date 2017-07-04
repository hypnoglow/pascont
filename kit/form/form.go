package form

import (
	"encoding/json"
	"fmt"
	"io"
)

// Form describes user input form.
type Form interface {
	// Validate validates the form and returns whether the Form is valid.
	Validate() bool

	// AddError reports an error occurred on form validation.
	AddError(message, field string, value interface{})

	// ValidationErrors returns the set of errors of the form.
	ValidationErrors() []FormError
}

// FormError is an error occurred on form validation.
type FormError struct {
	Message string
	Field   string
	Value   interface{}
}

// PopulateFormFromJSON populates the Form f from JSON body.
// It Closes the body.
func PopulateFormFromJSON(body io.ReadCloser, f Form) {
	d := json.NewDecoder(body)
	if err := d.Decode(f); err != nil {
		f.AddError(fmt.Sprintf("Request body is not a valid JSON"), "", nil)
	}

	body.Close()
}
