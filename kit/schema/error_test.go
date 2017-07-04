package schema

import (
	"reflect"
	"testing"
)

func TestNewError(t *testing.T) {
	cases := []struct {
		message       string
		field         string
		value         interface{}
		expectedError Error
	}{
		{
			message: "Error message",
			field:   "name",
			value:   "Igor",
			expectedError: Error{
				Message: "Error message",
				Field:   "name",
				Value:   "Igor",
			},
		},
	}

	for i, c := range cases {
		e := NewError(c.message, c.field, c.value)
		if !reflect.DeepEqual(e, c.expectedError) {
			t.Errorf(
				"testcase %d: Expected %#v but got %#v",
				i,
				c.expectedError,
				e,
			)
		}
	}
}

func TestErrorFromMessage(t *testing.T) {
	cases := []struct {
		message       string
		expectedError Error
	}{
		{
			message: "Error message",
			expectedError: Error{
				Message: "Error message",
				Field:   "",
				Value:   nil,
			},
		},
	}

	for i, c := range cases {
		e := ErrorFromMessage(c.message)
		if !reflect.DeepEqual(e, c.expectedError) {
			t.Errorf(
				"testcase %d: Expected %#v but got %#v",
				i,
				c.expectedError,
				e,
			)
		}
	}
}

func TestErrorsFromError(t *testing.T) {
	cases := []struct {
		error          Error
		expectedErrors []Error
	}{
		{
			error: Error{
				Message: "Error message",
				Field:   "name",
				Value:   "Igor",
			},
			expectedErrors: []Error{
				{
					Message: "Error message",
					Field:   "name",
					Value:   "Igor",
				},
			},
		},
	}

	for i, c := range cases {
		e := ErrorsFromError(c.error)
		if !reflect.DeepEqual(e, c.expectedErrors) {
			t.Errorf(
				"testcase %d: Expected %#v but got %#v",
				i,
				c.expectedErrors,
				e,
			)
		}
	}
}
