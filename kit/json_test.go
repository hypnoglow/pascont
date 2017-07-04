package kit

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/hypnoglow/pascont/kit/form"
	"github.com/hypnoglow/pascont/kit/schema"
)

func TestRespond(t *testing.T) {
	cases := []struct {
		statusCode   int
		body         []byte
		expectedCode int
		expectedBody []byte
	}{
		{
			statusCode:   200,
			body:         nil,
			expectedCode: 200,
			expectedBody: nil,
		},
		{
			statusCode:   400,
			body:         []byte(`{"a":1}`),
			expectedCode: 400,
			expectedBody: []byte(`{"a":1}`),
		},
		// TODO: test panic on body write error.
	}

	for i, c := range cases {
		r := httptest.NewRecorder()
		Respond(r, c.statusCode, c.body)

		if r.Result().StatusCode != c.expectedCode {
			t.Errorf(
				"testcase %d: Expected code to be %v but got %v",
				i,
				c.expectedCode,
				r.Result().StatusCode,
			)
		}

		body, err := ioutil.ReadAll(r.Result().Body)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(body, c.expectedBody) {
			t.Errorf(
				"testcase %d: Expected body to be %v but got %v",
				i,
				c.expectedBody,
				body,
			)
		}
	}
}

func TestRespondJSON(t *testing.T) {
	cases := []struct {
		statusCode     int
		payload        interface{}
		expectedCode   int
		expectedHeader http.Header
		expectedBody   string
	}{
		// If body is not nil, Respond should set Content-Type header.
		{
			// in
			statusCode: 200,
			payload: struct {
				Field string `json:"field"`
			}{
				Field: "value",
			},
			// out
			expectedCode:   200,
			expectedHeader: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
			expectedBody:   `{"field":"value"}`,
		},
		// If body is nil, Respond should not set Content-Type header.
		{
			// in
			statusCode: 201,
			payload:    nil,
			// out
			expectedCode:   201,
			expectedHeader: http.Header{},
			expectedBody:   "",
		},
	}

	for i, c := range cases {
		r := httptest.NewRecorder()
		RespondJSON(r, c.statusCode, c.payload)

		if r.Code != c.expectedCode {
			t.Errorf(
				"testcase %d: Expected status code to be %v but got %v\n",
				i,
				c.expectedCode,
				r.Code,
			)
		}

		if !reflect.DeepEqual(r.Result().Header, c.expectedHeader) {
			t.Errorf(
				"testcase %d: Expected header to be %v but got %v\n",
				i,
				c.expectedHeader,
				r.Result().Header,
			)
		}

		if r.Body.String() != c.expectedBody {
			t.Errorf(
				"testcase %d: Expected body to be %v, but got %v\n",
				i,
				c.expectedBody,
				r.Body.String(),
			)
		}
	}
}

func TestRespondWithFormErrors(t *testing.T) {
	cases := []struct {
		errors       []form.FormError
		expectedBody []byte
	}{
		{
			errors: []form.FormError{
				{
					Message: "A form error occured",
				},
			},
			expectedBody: []byte(`{"errors":[{"message":"A form error occured"}]}`),
		},
	}

	for i, c := range cases {
		r := httptest.NewRecorder()
		RespondWithFormErrors(r, 400, c.errors)

		b, err := ioutil.ReadAll(r.Result().Body)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(b, c.expectedBody) {
			t.Errorf(
				"testcase %d: Expected body to be %v but got %v",
				i,
				c.expectedBody,
				b,
			)
		}
	}
}

func TestRespondWithError(t *testing.T) {
	cases := []struct {
		error        schema.Error
		expectedBody []byte
	}{
		{
			error: schema.Error{
				Message: "An error occured",
			},
			expectedBody: []byte(`{"errors":[{"message":"An error occured"}]}`),
		},
	}

	for i, c := range cases {
		r := httptest.NewRecorder()
		RespondWithError(r, 400, c.error)

		b, err := ioutil.ReadAll(r.Result().Body)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(b, c.expectedBody) {
			t.Errorf(
				"testcase %d: Expected body to be %v but got %v",
				i,
				c.expectedBody,
				b,
			)
		}
	}
}
