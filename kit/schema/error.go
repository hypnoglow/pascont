package schema

type Error struct {
	Message string      `json:"message"`
	Field   string      `json:"field,omitempty"`
	Value   interface{} `json:"value,omitempty"`
}

func NewError(message, field string, value interface{}) Error {
	return Error{
		Message: message,
		Field:   field,
		Value:   value,
	}
}

func ErrorFromMessage(message string) Error {
	return NewError(message, "", nil)
}

func ErrorsFromError(err Error) []Error {
	return []Error{err}
}
