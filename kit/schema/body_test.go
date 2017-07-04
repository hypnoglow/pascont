package schema

import (
	"reflect"
	"testing"
)

func TestNewResultBody(t *testing.T) {
	result := struct {
		Field string `json:"field"`
	}{
		Field: "value",
	}
	meta := struct {
		Count int64 `json:"count"`
	}{
		Count: 1,
	}

	cases := []struct {
		result       interface{}
		meta         interface{}
		expectedBody *ResultBody
	}{
		{
			result: result,
			meta:   meta,
			expectedBody: &ResultBody{
				Result: result,
				Meta:   meta,
			},
		},
	}

	for i, c := range cases {
		actualBody := NewResultBody(c.result, c.meta)

		if !reflect.DeepEqual(actualBody, c.expectedBody) {
			t.Errorf(
				"testcase %d: Expected %v but got %v\n",
				i,
				c.expectedBody,
				actualBody,
			)
		}
	}
}

func TestNewResultsBody(t *testing.T) {
	result := struct {
		Field string `json:"field"`
	}{
		Field: "value",
	}
	meta := struct {
		Count int64 `json:"count"`
	}{
		Count: 1,
	}

	cases := []struct {
		results      []interface{}
		meta         interface{}
		expectedBody *ResultsBody
	}{
		{
			results: []interface{}{result},
			meta:    meta,
			expectedBody: &ResultsBody{
				Results: []interface{}{result},
				Meta:   meta,
			},
		},
	}

	for i, c := range cases {
		actualBody := NewResultsBody(c.results, c.meta)

		if !reflect.DeepEqual(actualBody, c.expectedBody) {
			t.Errorf(
				"testcase %d: Expected %v but got %v\n",
				i,
				c.expectedBody,
				actualBody,
			)
		}
	}
}

func TestNewErrorsBody(t *testing.T) {
	cases := []struct {
		errors       []Error
		expectedBody *ErrorsBody
	}{
		{
			errors: []Error{{Message: "test error"}},
			expectedBody: &ErrorsBody{
				Errors: []Error{
					{Message: "test error"},
				},
			},
		},
		{
			errors:       nil,
			expectedBody: nil,
		},
	}

	for i, c := range cases {
		actual := NewErrorsBody(c.errors)

		if !reflect.DeepEqual(actual, c.expectedBody) {
			t.Errorf(
				"testcase %d: Expected %v but got %v\n",
				i,
				c.expectedBody,
				actual,
			)
		}
	}
}

func TestErrorsBodyFromError(t *testing.T) {
	cases := []struct {
		error        Error
		expectedBody *ErrorsBody
	}{
		{
			error: Error{Message: "Test"},
			expectedBody: &ErrorsBody{
				Errors: []Error{
					{Message: "Test"},
				},
			},
		},
	}

	for i, c := range cases {
		actualBody := ErrorsBodyFromError(c.error)
		if !reflect.DeepEqual(actualBody, c.expectedBody) {
			t.Errorf(
				"testcase %d: Expected %v but got %v\n",
				i,
				c.expectedBody,
				actualBody,
			)
		}
	}
}
