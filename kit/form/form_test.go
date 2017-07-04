package form

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

type SampleTestForm struct {
	BaseForm
	Name string `json:"name"`
}

func (s *SampleTestForm) Validate() bool {
	return true
}

func TestPopulateFormFromJSON(t *testing.T) {
	cases := []struct {
		body         io.ReadCloser
		form         Form
		expectedForm Form
	}{
		{
			body: ioutil.NopCloser(bytes.NewBufferString(`{"name":"Igor"}`)),
			form: &SampleTestForm{},
			expectedForm: &SampleTestForm{
				Name: "Igor",
			},
		},
		{
			body: ioutil.NopCloser(bytes.NewBufferString(`{"wrong json"}`)),
			form: &SampleTestForm{},
			expectedForm: func() Form {
				form := &SampleTestForm{}
				form.AddError("Request body is not a valid JSON", "", nil)
				return form
			}(),
		},
	}

	for i, c := range cases {
		PopulateFormFromJSON(c.body, c.form)

		if !reflect.DeepEqual(c.form, c.expectedForm) {
			t.Errorf(
				"testcase %d: Expected form to be %v but got %v",
				i,
				c.expectedForm,
				c.form,
			)
		}
	}
}
