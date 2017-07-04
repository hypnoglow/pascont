package kit

import (
	"bytes"
	"testing"
)

func TestNewJSONBody(t *testing.T) {
	cases := []struct {
		payload       interface{}
		expectedBody  []byte
		expectedPanic bool
	}{
		{
			payload: struct {
				Message string `json:"message"`
			}{
				Message: "test",
			},
			expectedBody:  []byte(`{"message":"test"}`),
			expectedPanic: false,
		},
		{
			payload:       nil,
			expectedBody:  nil,
			expectedPanic: false,
		},
		{
			payload:       make(chan int),
			expectedPanic: true,
		},
	}

	for i, c := range cases {
		func() {
			defer func() {
				if err := recover(); err != nil {
					if !c.expectedPanic {
						t.Errorf("testcase %d: Unexpected panic %v", err)
					}
				}
			}()

			actual := NewJSONBody(c.payload)

			if !bytes.Equal(actual, c.expectedBody) {
				t.Errorf(
					"testcase %d: Expected %v but got %v\n",
					i,
					c.expectedBody,
					actual,
				)
			}
		}()
	}
}
