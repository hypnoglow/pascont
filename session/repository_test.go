package session

import "testing"

func TestSessionRepositoryError_Error(t *testing.T) {
	msg := "some error"

	err := repositoryError(msg)
	if err.Error() != msg {
		t.Errorf("Expected %v but got %v'n", msg, err.Error())
	}
}
