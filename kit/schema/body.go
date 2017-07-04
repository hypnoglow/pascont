package schema

type ResultBody struct {
	Result interface{} `json:"result"`
	Meta   interface{} `json:"meta,omitempty"`
}

func NewResultBody(result interface{}, meta interface{}) (body *ResultBody) {
	return &ResultBody{result, meta}
}

type ResultsBody struct {
	Results interface{} `json:"results"`
	Meta    interface{} `json:"meta,omitempty"`
}

func NewResultsBody(results interface{}, meta interface{}) (body *ResultsBody) {
	return &ResultsBody{results, meta}
}

type ErrorsBody struct {
	Errors []Error `json:"errors"`
}

func NewErrorsBody(errors []Error) (body *ErrorsBody) {
	if errors == nil {
		return nil
	}

	body = &ErrorsBody{}
	for _, e := range errors {
		body.Errors = append(body.Errors, e)
	}

	return body
}

func ErrorsBodyFromError(error Error) (body *ErrorsBody) {
	return NewErrorsBody([]Error{error})
}
