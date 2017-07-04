package kit

import (
	"net/http"

	"github.com/hypnoglow/pascont/kit/form"
	"github.com/hypnoglow/pascont/kit/schema"
	"github.com/pkg/errors"
)

// Respond writes status code and body to w.
func Respond(w http.ResponseWriter, status int, body []byte) {
	w.WriteHeader(status)
	_, err := w.Write(body)
	if err != nil {
		panic(errors.Wrap(err, "Failed to write body to ResponseWriter"))
	}
}

// RespondJSON writes status code and json encoded payload to w.
// If payload is not nil, header Content-Type is set as application/json.
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	if payload != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	Respond(w, status, NewJSONBody(payload))
}

// RespondWithFormErrors is a shortcut method for responding with
// form errors as schema.ErrorsBody
func RespondWithFormErrors(w http.ResponseWriter, status int, errors []form.FormError) {
	RespondJSON(w, status, schema.NewErrorsBody(
		newErrorsFromFormErrors(errors),
	))
}

// RespondWithError is a shortcut method for responding with
// schema.Error as schema.ErrorsBody
func RespondWithError(w http.ResponseWriter, status int, error schema.Error) {
	RespondJSON(w, status, schema.ErrorsBodyFromError(error))
}

func newErrorsFromFormErrors(fe []form.FormError) (errors []schema.Error) {
	errors = make([]schema.Error, len(fe))
	for i, e := range fe {
		errors[i] = schema.Error(e)
	}
	return errors
}
