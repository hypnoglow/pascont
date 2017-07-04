package form

import (
	"testing"
	"reflect"
)

func TestBaseForm_AddError(t *testing.T) {
	err := FormError{Message: "Test error", Field: "field", Value: "value"}
	expected := []FormError{err}

	f := BaseForm{}
	f.AddError(err.Message, err.Field, err.Value)

	if !reflect.DeepEqual(expected, f.errors) {
		t.Errorf("Expected %v but got %v\n", expected, f.errors)
	}
}

func TestBaseForm_ValidationErrors(t *testing.T) {
	expected := []FormError{{Message: "Test error", Field: "field", Value: "value"}}
	f := BaseForm{errors: expected}

	if !reflect.DeepEqual(expected, f.ValidationErrors()) {
		t.Errorf("Expected %v but got %v\n", expected, f.ValidationErrors())
	}
}