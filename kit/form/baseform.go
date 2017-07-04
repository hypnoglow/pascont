package form

// BaseForm is a base form supporting errors.
type BaseForm struct {
	errors []FormError
}

func (f *BaseForm) AddError(message, field string, value interface{}) {
	f.errors = append(f.errors, FormError{message, field, value})
}

func (f BaseForm) ValidationErrors() []FormError {
	return f.errors
}
