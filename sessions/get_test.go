package sessions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/hypnoglow/pascont/kit/middleware"
	"github.com/hypnoglow/pascont/session"
)

func TestRestController_GetSession(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	later := now.Add(time.Hour * 24)
	fakeLogger := log.New(ioutil.Discard, "", log.LstdFlags)

	cases := []struct {
		caseName string
		// in
		sessRepo                 session.Repository
		reqContextSessionIDValue interface{}
		reqContextPathIDValue    interface{}
		reqAuthHeader            string
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
			// out
			expectedCode: http.StatusInternalServerError,
			expectedBody: bytes.NewBuffer(nil),
		},
		{
			caseName: "Successful",
			sessRepo: session.NewFakeRepository(
				nil,
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
			// out
			expectedCode: http.StatusOK,
			expectedBody: bytes.NewBufferString(fmt.Sprintf(`{
				"result": {
					"token":"",
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
			nil,
			nil,
			nil,
			nil,
			Options{},
		)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, PathSession, nil)
		req = req.WithContext(context.WithValue(req.Context(), middleware.ContextKeyPathID{}, c.reqContextPathIDValue))
		req = req.WithContext(context.WithValue(req.Context(), middleware.ContextKeySessionID{}, c.reqContextSessionIDValue))
		ctrl.GetSession(w, req)

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
