package accounts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/hypnoglow/pascont/account"
	"github.com/hypnoglow/pascont/hasher"
	"github.com/hypnoglow/pascont/kit/form"
)

func TestPostAccountForm_Validate(t *testing.T) {
	cases := []struct {
		caseName       string
		Name           string
		Password       string
		expectedErrors []form.FormError
	}{
		{
			caseName: "Name is too short",
			Name:     "i",
			Password: "password",
			expectedErrors: []form.FormError{
				{
					Message: fmt.Sprintf("Name must be at least %d bytes", nameMinLen),
					Field:   "name",
					Value:   "i",
				},
			},
		},
		{
			caseName: "Name is too long",
			Name:     "1234567890123456789012345678901234567890123456789012345678901234567890",
			Password: "password",
			expectedErrors: []form.FormError{
				{
					Message: fmt.Sprintf("Name must be at most %d bytes", nameMaxLen),
					Field:   "name",
					Value:   "1234567890123456789012345678901234567890123456789012345678901234567890",
				},
			},
		},
		{
			caseName: "Password is too short",
			Name:     "email@email.com",
			Password: "pass",
			expectedErrors: []form.FormError{
				{
					Message: fmt.Sprintf("Password must be at least %d bytes", passwordMinLen),
					Field:   "password",
					Value:   nil,
				},
			},
		},
	}

	for i, c := range cases {
		f := &postAccountForm{
			Name:     c.Name,
			Password: c.Password,
		}

		f.Validate()

		if !reflect.DeepEqual(f.ValidationErrors(), c.expectedErrors) {
			t.Errorf(
				fmt.Sprintf(
					"testcase %d %s\nExpected validation errors to equal\n%#v\n, but got\n%#v\n",
					i,
					c.caseName,
					c.expectedErrors,
					f.ValidationErrors(),
				),
			)
		}
	}
}

func TestRestController_PostAccounts(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	fakeLogger := log.New(ioutil.Discard, "", log.LstdFlags)

	cases := []struct {
		caseName          string
		accRepo           account.Repository
		hasher            hasher.Hasher
		opts              Options
		form              postAccountForm
		reqBody           io.Reader
		expectedCode      int
		expectedHeaderMap http.Header
		expectedBody      *bytes.Buffer
	}{
		{
			caseName: "Too short account name should result in 400",
			accRepo:  account.NewFakeRepository(nil, nil, nil, nil),
			reqBody: bytes.NewBufferString(`{
				"name":"i",
				"password":"password"
			}`),
			expectedCode:      http.StatusBadRequest,
			expectedHeaderMap: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
			expectedBody: bytes.NewBufferString(`{
				"errors":[
					{
						"message":"Name must be at least 4 bytes",
						"field":"name",
						"value":"i"
					}
				]
			}`),
		},
		{
			caseName: "Error on repository Exists",
			accRepo: account.NewFakeRepository(
				nil,
				nil,
				[]account.FakeRepositoryExistsResult{
					{
						Exists: false,
						Error:  fmt.Errorf("Exists failed"),
					},
				},
				nil,
			),
			reqBody: bytes.NewBufferString(`{
				"name":"email@email.com",
				"password":"password"
			}`),
			expectedCode:      http.StatusInternalServerError,
			expectedHeaderMap: http.Header{},
			expectedBody:      bytes.NewBuffer(nil),
		},
		{
			caseName: "Account with such name already exists",
			accRepo: account.NewFakeRepository(
				nil,
				nil,
				[]account.FakeRepositoryExistsResult{
					{
						Exists: true,
						Error:  nil,
					},
				},
				nil,
			),
			reqBody: bytes.NewBufferString(`{
				"name":"email@email.com",
				"password":"password"
			}`),
			expectedCode:      http.StatusConflict,
			expectedHeaderMap: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
			expectedBody: bytes.NewBufferString(`{
				"errors":[
					{
						"message":"Account with name email@email.com already exists",
						"field":"name",
						"value":"email@email.com"
					}
				]
			}`),
		},
		{
			caseName: "Error on hasher.GenerateHashFromPassword should result in 500",
			accRepo: account.NewFakeRepository(
				nil,
				nil,
				[]account.FakeRepositoryExistsResult{
					{
						Exists: false,
						Error:  nil,
					},
				},
				nil,
			),
			hasher: hasher.NewFakeHasher(
				[]hasher.FakeGenerateHashFromPasswordResult{
					{
						Hash:  nil,
						Error: fmt.Errorf("Failed to generate hash"),
					},
				},
				nil,
			),
			reqBody: bytes.NewBufferString(`{
				"name":"email@email.com",
				"password":"password"
			}`),
			expectedCode:      http.StatusInternalServerError,
			expectedHeaderMap: http.Header{},
			expectedBody:      bytes.NewBuffer(nil),
		},
		{
			caseName: "Error on repository Save",
			accRepo: account.NewFakeRepository(
				[]account.FakeRepositoryAcceptResult{
					{
						Account: nil,
						Error:   fmt.Errorf("Accept failed"),
					},
				},
				nil,
				[]account.FakeRepositoryExistsResult{
					{
						Exists: false,
						Error:  nil,
					},
				},
				nil,
			),
			hasher: hasher.NewFakeHasher(
				[]hasher.FakeGenerateHashFromPasswordResult{
					{
						Hash:  []byte("password_hash"),
						Error: nil,
					},
				},
				nil,
			),
			reqBody: bytes.NewBufferString(`{
				"name":"email@email.com",
				"password":"password"
			}`),
			expectedCode:      http.StatusInternalServerError,
			expectedHeaderMap: http.Header{},
			expectedBody:      bytes.NewBuffer(nil),
		},
		{
			caseName: "Successful",
			accRepo: account.NewFakeRepository(
				[]account.FakeRepositoryAcceptResult{
					{
						Account: &account.Account{
							ID:        123,
							Name:      "email@email.com",
							CreatedAt: now,
							UpdatedAt: now,
						},
						Error: nil,
					},
				},
				nil,
				[]account.FakeRepositoryExistsResult{
					{
						Exists: false,
						Error:  nil,
					},
				},
				nil,
			),
			hasher: hasher.NewFakeHasher(
				[]hasher.FakeGenerateHashFromPasswordResult{
					{
						Hash:  []byte("password_hash"),
						Error: nil,
					},
				},
				nil,
			),
			reqBody: bytes.NewBufferString(`{
				"name":"email@email.com",
				"password":"password"
			}`),
			expectedCode:      http.StatusCreated,
			expectedHeaderMap: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
			expectedBody: bytes.NewBufferString(fmt.Sprintf(`{
				"result":{
					"id":123,
					"name":"email@email.com",
					"createdAt":"%s",
					"updatedAt":"%s"
				}
			}`, now.Format(time.RFC3339), now.Format(time.RFC3339))),
		},
	}

	for i, c := range cases {
		ctrl := NewRestController(
			fakeLogger,
			c.accRepo,
			c.hasher,
			c.opts,
		)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, PathAccounts, c.reqBody)

		ctrl.PostAccounts(w, req)

		if w.Code != c.expectedCode {
			t.Errorf(
				"testcase %d %s:\nExpected status code to be %v, but got %v\n",
				i,
				c.caseName,
				c.expectedCode,
				w.Code,
			)
		}

		if !reflect.DeepEqual(w.HeaderMap, c.expectedHeaderMap) {
			t.Errorf(
				"testcase %d %s:\nExpected header map to be %v, but got %v\n",
				i,
				c.caseName,
				c.expectedHeaderMap,
				w.HeaderMap,
			)
		}

		// To make possible pretty formatting in test cases, we use json.Compact:
		b := c.expectedBody.Bytes()
		if b != nil {
			c.expectedBody.Reset()
			err := json.Compact(c.expectedBody, b)
			if err != nil {
				t.Errorf(
					"testcase %d %s:\nFailed to compact JSON body: %s",
					i,
					c.caseName,
					err,
				)
			}
		}

		if !bytes.Equal(w.Body.Bytes(), c.expectedBody.Bytes()) {
			t.Errorf(
				"testcase %d %s:\nExpected body to be\n%v\nbut got\n%v\n",
				i,
				c.caseName,
				c.expectedBody,
				w.Body,
			)
		}
	}
}
