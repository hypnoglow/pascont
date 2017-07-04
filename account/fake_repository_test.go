package account

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNewFakeRepository(t *testing.T) {
	NewFakeRepository(nil, nil, nil, nil)
}

func TestFakeRepository_Accept(t *testing.T) {
	fakeAcc := &Account{123, "email@email.com", time.Now(), time.Now()}

	cases := []struct {
		acceptResults    []FakeRepositoryAcceptResult
		expectedAccounts []*Account
		expectedErrors   []error
	}{
		{
			acceptResults: []FakeRepositoryAcceptResult{
				{
					Account: fakeAcc,
					Error:   nil,
				},
				{
					Account: nil,
					Error:   fmt.Errorf("Some error"),
				},
			},
			expectedAccounts: []*Account{
				fakeAcc,
				nil,
			},
			expectedErrors: []error{
				nil,
				fmt.Errorf("Some error"),
			},
		},
	}

	for i, c := range cases {
		fr := fakeRepository{
			acceptResults:       c.acceptResults,
			acceptResultCounter: 0,
		}

		for resultIndex := range c.acceptResults {
			acc, err := fr.Accept(Application{})

			if !reflect.DeepEqual(acc, c.expectedAccounts[resultIndex]) {
				t.Errorf(
					"testcase %d: Expected Account to be %v but got %v",
					i,
					c.expectedAccounts[resultIndex],
					acc,
				)
			}

			if !reflect.DeepEqual(err, c.expectedErrors[resultIndex]) {
				t.Errorf(
					"testcase %d: Expected Error to be %v but got %v",
					i,
					c.expectedErrors[resultIndex],
					err,
				)
			}
		}

	}
}

func TestFakeRepository_Save(t *testing.T) {
	// TODO
}

func TestFakeRepository_FindWithPasswordHashByUsername(t *testing.T) {
	// TODO
}

func TestFakeRepository_Exists(t *testing.T) {
	// TODO
}
