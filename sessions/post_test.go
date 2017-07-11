package sessions

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
	"github.com/hypnoglow/pascont/notary"
	"github.com/hypnoglow/pascont/packer"
	"github.com/hypnoglow/pascont/session"
)

func TestPostSessionForm_Validate(t *testing.T) {
	cases := []struct {
		caseName       string
		Name           string
		Password       string
		expectedErrors []form.FormError
	}{
		{
			caseName: "Name is empty",
			Name:     "",
			Password: "password",
			expectedErrors: []form.FormError{
				{
					Message: "Name must not be empty",
					Field:   "name",
					Value:   "",
				},
			},
		},
		{
			caseName: "PasswordHash is empty",
			Name:     "email@email.com",
			Password: "",
			expectedErrors: []form.FormError{
				{
					Message: "PasswordHash must not be empty",
					Field:   "password",
					Value:   "",
				},
			},
		},
	}

	for i, c := range cases {
		f := &postSessionForm{
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

func TestRestController_PostSessions(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	fakeLogger := log.New(ioutil.Discard, "", log.LstdFlags)
	uuidProducer := func() string {
		return "12345678-90ab-cdef-0123-4567890abcde"
	}

	cases := []struct {
		caseName          string
		sessRepo          session.Repository
		accRepo           account.Repository
		notary            notary.Notary
		packer            packer.Packer
		hasher            hasher.Hasher
		reqBody           io.Reader
		expectedCode      int
		expectedHeaderMap http.Header
		expectedBody      *bytes.Buffer
	}{
		{
			caseName: "Empty password should result in 400",
			reqBody: bytes.NewBufferString(`{
				"name":"email@email.com",
				"password":""
			}`),
			expectedCode:      http.StatusBadRequest,
			expectedHeaderMap: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
			expectedBody: bytes.NewBufferString(`{
				"errors":[
					{
						"message":"PasswordHash must not be empty",
						"field":"password",
						"value":""
					}
				]
			}`),
		},
		{
			caseName: "Non-existent name should result in 401",
			accRepo: account.NewFakeRepository(
				nil,
				nil,
				nil,
				[]account.FakeRepositoryFindWithPasswordHashByUsernameResult{
					{
						Account: nil,
						Error:   account.ErrNotFound,
					},
				},
			),
			reqBody: bytes.NewBufferString(`{
				"name":"nonexistent@email.com",
				"password":"password"
			}`),
			expectedCode:      http.StatusUnauthorized,
			expectedHeaderMap: http.Header{},
			expectedBody:      bytes.NewBuffer(nil),
		},
		{
			caseName: "Error on FindWithPasswordHashByUsername should result in 500",
			accRepo: account.NewFakeRepository(
				nil,
				nil,
				nil,
				[]account.FakeRepositoryFindWithPasswordHashByUsernameResult{
					{
						Account: nil,
						Error:   fmt.Errorf("FindWithPasswordHashByUsername failed"),
					},
				},
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
			caseName: "Error on checking password should result in 401",
			accRepo: account.NewFakeRepository(
				nil,
				nil,
				nil,
				[]account.FakeRepositoryFindWithPasswordHashByUsernameResult{
					{
						Account: &account.Account{
							ID:        123,
							Name:      "email@email.com",
							CreatedAt: now,
							UpdatedAt: now,
						},
						PasswordHash: []byte("password_hash"),
						Error:        nil,
					},
				},
			),
			hasher: hasher.NewFakeHasher(
				nil,
				[]hasher.FakeCompareHashWithPasswordResult{
					{
						Error: fmt.Errorf("No match"),
					},
				},
			),
			reqBody: bytes.NewBufferString(`{
				"name":"email@email.com",
				"password":"password"
			}`),
			expectedCode:      http.StatusUnauthorized,
			expectedHeaderMap: http.Header{},
			expectedBody:      bytes.NewBuffer(nil),
		},
		{
			caseName: "Error on sessionRepo.Save should result in 500",
			accRepo: account.NewFakeRepository(
				nil,
				[]account.FakeRepositorySaveResult{
					{
						Error: nil,
					},
				},
				nil,
				[]account.FakeRepositoryFindWithPasswordHashByUsernameResult{
					{
						Account: &account.Account{
							ID:        123,
							Name:      "email@email.com",
							CreatedAt: now,
							UpdatedAt: now,
						},
						PasswordHash: []byte("password_hash"),
						Error:        nil,
					},
				},
			),
			sessRepo: session.NewFakeRepository(
				[]session.FakeRepositorySaveResult{
					{
						Error: fmt.Errorf("Save failed"),
					},
				},
				nil,
			),
			hasher: hasher.NewFakeHasher(
				nil,
				[]hasher.FakeCompareHashWithPasswordResult{
					{
						Error: nil,
					},
				},
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
			caseName: "Error on sess.Token should result in 500",
			sessRepo: session.NewFakeRepository(
				[]session.FakeRepositorySaveResult{
					{
						Error: nil,
					},
				},
				nil,
			),
			accRepo: account.NewFakeRepository(
				nil,
				nil,
				nil,
				[]account.FakeRepositoryFindWithPasswordHashByUsernameResult{
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
			),
			notary: notary.NewFakeNotary(
				[]notary.FakeNotarySignResult{
					{
						Signature: []byte("signature"),
					},
				},
				nil,
			),
			packer: packer.NewFakePacker(
				[]packer.FakePackerPackResult{
					{
						Pack:  nil,
						Error: fmt.Errorf("Pack failed"),
					},
				},
				nil,
			),
			hasher: hasher.NewFakeHasher(
				nil,
				[]hasher.FakeCompareHashWithPasswordResult{
					{
						Error: nil,
					},
				},
			),
			reqBody: bytes.NewBufferString(
				`{"name":"email@email.com","password":"password"}`,
			),
			expectedCode:      http.StatusInternalServerError,
			expectedHeaderMap: http.Header{},
			expectedBody:      bytes.NewBuffer(nil),
		},
		{
			caseName: "Successful",
			sessRepo: session.NewFakeRepository(
				[]session.FakeRepositorySaveResult{
					{
						Error: nil,
					},
				},
				nil,
			),
			accRepo: account.NewFakeRepository(
				nil,
				nil,
				nil,
				[]account.FakeRepositoryFindWithPasswordHashByUsernameResult{
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
			),
			notary: notary.NewFakeNotary(
				[]notary.FakeNotarySignResult{
					{
						Signature: []byte("signature"),
					},
				},
				nil,
			),
			packer: packer.NewFakePacker(
				[]packer.FakePackerPackResult{
					{
						Pack:  []byte("pack"),
						Error: nil,
					},
				},
				nil,
			),
			hasher: hasher.NewFakeHasher(
				nil,
				[]hasher.FakeCompareHashWithPasswordResult{
					{
						Error: nil,
					},
				},
			),
			reqBody: bytes.NewBufferString(
				`{"name":"email@email.com","password":"password"}`,
			),
			expectedCode: http.StatusCreated,
			expectedHeaderMap: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			expectedBody: bytes.NewBufferString(fmt.Sprintf(`{
				"result": {
					"token":"pack",
					"id":"12345678-90ab-cdef-0123-4567890abcde",
					"accountID":123,
					"createdAt":"%s",
					"expiresAt":"%s"
				}
			}`, now.Format(time.RFC3339), now.Add(sessionDefaultDuration).Format(time.RFC3339))),
		},
	}

	for i, c := range cases {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("testcase %d %s\npanic: %+v", i, c.caseName, err)
			}
		}()

		ctrl := NewRestController(
			fakeLogger,
			c.accRepo,
			c.sessRepo,
			c.notary,
			c.packer,
			c.hasher,
			uuidProducer,
			Options{},
		)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, PathSessions, c.reqBody)

		ctrl.PostSessions(w, req)

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
