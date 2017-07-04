package sessions

import (
	"bytes"
	"context"
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

	"github.com/hypnoglow/pascont/kit/form"
	"github.com/hypnoglow/pascont/kit/middleware"
	"github.com/hypnoglow/pascont/notary"
	"github.com/hypnoglow/pascont/packer"
	"github.com/hypnoglow/pascont/session"
)

func TestPatchSessionForm_Validate(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	cases := []struct {
		caseName       string
		expiresAt      time.Time
		expectedErrors []form.FormError
	}{
		{
			caseName:  "ExpiresAt is earlier than now",
			expiresAt: now.Add(-time.Hour * 24),
			expectedErrors: []form.FormError{
				{
					Message: fmt.Sprintf(
						"ExpiresAt must be greater than %s",
						now.Format(time.RFC3339),
					),
					Field: "expiresAt",
					Value: now.Add(-time.Hour * 24),
				},
			},
		},
		{
			caseName:  "ExpiresAt is too much later than now",
			expiresAt: now.Add(sessionMaxDuration + time.Hour*24),
			expectedErrors: []form.FormError{
				{
					Message: fmt.Sprintf(
						"ExpiresAt must be not greater than %s",
						now.Add(sessionMaxDuration).Format(time.RFC3339),
					),
					Field: "expiresAt",
					Value: now.Add(sessionMaxDuration + time.Hour*24),
				},
			},
		},
	}

	for i, c := range cases {
		f := &patchSessionForm{
			ExpiresAt: c.expiresAt,
		}

		f.Validate()

		if !reflect.DeepEqual(f.ValidationErrors(), c.expectedErrors) {
			t.Errorf(
				fmt.Sprintf(
					"testcase %d %s\nExpected validation errors to equal\n%v\n, but got\n%v\n",
					i,
					c.caseName,
					c.expectedErrors,
					f.ValidationErrors(),
				),
			)
		}
	}
}

func TestRestController_PatchSession(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	earlier := now.Add(-time.Hour * 24)
	later := now.Add(time.Hour * 24)
	fakeLogger := log.New(ioutil.Discard, "", log.LstdFlags)

	cases := []struct {
		caseName string
		// in
		sessRepo                 session.Repository
		notary                   notary.Notary
		packer                   packer.Packer
		reqContextSessionIDValue interface{}
		reqContextPathIDValue    interface{}
		reqBody                  io.Reader
		// out
		expectedCode      int
		expectedHeaderMap http.Header
		expectedBody      *bytes.Buffer
	}{
		{
			caseName:              "The id from request context is not a string",
			reqContextPathIDValue: 123,
			// out
			expectedCode: http.StatusNotFound,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName:              "There is no session ID from Authorization Header",
			reqContextPathIDValue: "99999999-90ab-cdef-0123-4567890abcde",
			// out
			expectedCode: http.StatusUnauthorized,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName:                 "Session ids from token and from url path do not match",
			reqContextSessionIDValue: "12345678-90ab-cdef-0123-4567890abcde",
			reqContextPathIDValue:    "99999999-90ab-cdef-0123-4567890abcde",
			// out
			expectedCode: http.StatusUnauthorized,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName:                 "Invalid input form",
			reqContextSessionIDValue: "12345678-90ab-cdef-0123-4567890abcde",
			reqContextPathIDValue:    "12345678-90ab-cdef-0123-4567890abcde",
			reqBody: bytes.NewBufferString(fmt.Sprintf(`{
				"expiresAt":"%s"
			}`, earlier.Format(time.RFC3339))),
			// out
			expectedCode: http.StatusBadRequest,
			expectedHeaderMap: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			expectedBody: bytes.NewBufferString(fmt.Sprintf(`{
				"errors": [
					{
						"message":"ExpiresAt must be greater than %s",
						"field":"expiresAt",
						"value":"%s"
					}
				]
			}`, now.Format(time.RFC3339), earlier.Format(time.RFC3339))),
		},
		{
			caseName: "sessionRepo.FindByID returns ErrNotFound",
			sessRepo: session.NewFakeRepository(
				nil,
				[]session.FakeRepositoryFindByIDResult{
					{
						Error: session.ErrNotFound,
					},
				},
			),
			reqContextSessionIDValue: "12345678-90ab-cdef-0123-4567890abcde",
			reqContextPathIDValue:    "12345678-90ab-cdef-0123-4567890abcde",
			reqBody: bytes.NewBufferString(fmt.Sprintf(`{
				"expiresAt":"%s"
			}`, later.Format(time.RFC3339))),
			// out
			expectedCode: http.StatusNotFound,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName: "sessionRepo.FindByID returns ErrExpired",
			sessRepo: session.NewFakeRepository(
				nil,
				[]session.FakeRepositoryFindByIDResult{
					{
						Error: session.ErrExpired,
					},
				},
			),
			reqContextSessionIDValue: "12345678-90ab-cdef-0123-4567890abcde",
			reqContextPathIDValue:    "12345678-90ab-cdef-0123-4567890abcde",
			reqBody: bytes.NewBufferString(fmt.Sprintf(`{
				"expiresAt":"%s"
			}`, later.Format(time.RFC3339))),
			// out
			expectedCode: http.StatusNotFound,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName: "sessionRepo.FindByID failed",
			sessRepo: session.NewFakeRepository(
				nil,
				[]session.FakeRepositoryFindByIDResult{
					{
						Error: fmt.Errorf("FindByID failed"),
					},
				},
			),
			reqContextSessionIDValue: "12345678-90ab-cdef-0123-4567890abcde",
			reqContextPathIDValue:    "12345678-90ab-cdef-0123-4567890abcde",
			reqBody: bytes.NewBufferString(fmt.Sprintf(`{
				"expiresAt":"%s"
			}`, later.Format(time.RFC3339))),
			// out
			expectedCode: http.StatusInternalServerError,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName: "sessionRepo.Save failed",
			sessRepo: session.NewFakeRepository(
				[]session.FakeRepositorySaveResult{
					{
						Error: fmt.Errorf("Save failed"),
					},
				},
				[]session.FakeRepositoryFindByIDResult{
					{
						Session: &session.Session{
							ID:        "12345678-90ab-cdef-0123-4567890abcde",
							AccountID: 132,
							CreatedAt: now,
							ExpiresAt: later,
						},
						Error: nil,
					},
				},
			),
			reqContextSessionIDValue: "12345678-90ab-cdef-0123-4567890abcde",
			reqContextPathIDValue:    "12345678-90ab-cdef-0123-4567890abcde",
			reqBody: bytes.NewBufferString(fmt.Sprintf(`{
				"expiresAt":"%s"
			}`, later.Format(time.RFC3339))),
			// out
			expectedCode: http.StatusInternalServerError,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName: "Issuing a new Token failed",
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
						Error: fmt.Errorf("Failed to Pack"),
					},
				},
				nil,
			),
			sessRepo: session.NewFakeRepository(
				[]session.FakeRepositorySaveResult{
					{
						Error: nil,
					},
				},
				[]session.FakeRepositoryFindByIDResult{
					{
						Session: &session.Session{
							ID:        "12345678-90ab-cdef-0123-4567890abcde",
							AccountID: 123,
							CreatedAt: now,
							ExpiresAt: later,
						},
						Error: nil,
					},
				},
			),
			reqContextSessionIDValue: "12345678-90ab-cdef-0123-4567890abcde",
			reqContextPathIDValue:    "12345678-90ab-cdef-0123-4567890abcde",
			reqBody: bytes.NewBufferString(fmt.Sprintf(`{
				"expiresAt":"%s"
			}`, later.Format(time.RFC3339))),
			// out
			expectedCode: http.StatusInternalServerError,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName: "Successful",
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
						Pack:  []byte("new-token"),
						Error: nil,
					},
				},
				nil,
			),
			sessRepo: session.NewFakeRepository(
				[]session.FakeRepositorySaveResult{
					{
						Error: nil,
					},
				},
				[]session.FakeRepositoryFindByIDResult{
					{
						Session: &session.Session{
							ID:        "12345678-90ab-cdef-0123-4567890abcde",
							AccountID: 123,
							CreatedAt: now,
							ExpiresAt: later,
						},
						Error: nil,
					},
				},
			),
			reqContextSessionIDValue: "12345678-90ab-cdef-0123-4567890abcde",
			reqContextPathIDValue:    "12345678-90ab-cdef-0123-4567890abcde",
			reqBody: bytes.NewBufferString(fmt.Sprintf(`{
				"expiresAt":"%s"
			}`, later.Format(time.RFC3339))),
			// out
			expectedCode: http.StatusOK,
			expectedBody: bytes.NewBufferString(fmt.Sprintf(`{
				"result": {
					"token":"new-token",
					"id":"12345678-90ab-cdef-0123-4567890abcde",
					"accountID":123,
					"createdAt":"%s",
					"expiresAt":"%s"
				}
			}`, now.Format(time.RFC3339), later.Format(time.RFC3339))),
		},
	}

	for i, c := range cases {
		ctrl := NewRestController(
			fakeLogger,
			nil,
			c.sessRepo,
			c.notary,
			c.packer,
			nil,
			nil,
			Options{},
		)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, PathSession, c.reqBody)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ContextKeyPathID{}, c.reqContextPathIDValue))
		req = req.WithContext(context.WithValue(req.Context(), middleware.ContextKeySessionID{}, c.reqContextSessionIDValue))
		ctrl.PatchSession(w, req)

		if w.Code != c.expectedCode {
			t.Errorf(
				"testcase %d %s:\nExpected status code to be %v, but got %v\n",
				i,
				c.caseName,
				c.expectedCode,
				w.Code,
			)

			if !reflect.DeepEqual(w.HeaderMap, c.expectedHeaderMap) {
				t.Errorf(
					"testcase %d %s:\nExpected header map to be %v, but got %v\n",
					i,
					c.caseName,
					c.expectedHeaderMap,
					w.HeaderMap,
				)
			}
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
