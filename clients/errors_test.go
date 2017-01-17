// Copyright 2017 Jeff Nickoloff "jeff@allingeek.com"
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clients

import (
	"errors"
	"net/http"
	"testing"
)

func TestRetriableError(t *testing.T) {
	var e ClientError
	e = RetriableError{E: nil}
	if !e.IsRetriable() {
		t.Fatal(`RetriableError was notretriable.`)
	}
	if e.Error() != nil {
		t.Fatal(`RetriableError with nil error returned non-nil error.`)
	}
	e = RetriableError{E: errors.New(`non-nil error`)}
	if !e.IsRetriable() {
		t.Fatal(`RetriableError was notretriable.`)
	}
	if e.Error() == nil {
		t.Fatal(`RetriableError with non-nil error returned nil error.`)
	}
}

func TestNonRetriableError(t *testing.T) {
	var e ClientError
	e = NonRetriableError{E: nil}
	if e.IsRetriable() {
		t.Fatal(`NonRetriableError was retriable.`)
	}
	if e.Error() != nil {
		t.Fatal(`NonRetriableError with nil error returned non-nil error.`)
	}
	e = NonRetriableError{E: errors.New(`non-nil error`)}
	if e.IsRetriable() {
		t.Fatal(`NonRetriableError was retriable.`)
	}
	if e.Error() == nil {
		t.Fatal(`NonRetriableError with non-nil error returned nil error.`)
	}
}

func TestWrapHttpResponseError(t *testing.T) {
	var e ClientError
	var ir, rr *http.Response

	rr, e = WrapHttpResponseError(nil, nil)
	if e != nil {
		t.Fatal(`Nil response and nil error input resulted in a non-nil err`)
	}
	if rr != nil {
		t.Fatal(`Nil response and nil error input resulted in a non-nil response`)
	}

	rr, e = WrapHttpResponseError(nil, errors.New(`Junk error`))
	if rr != nil {
		t.Fatal(`Nil response and non-nil error resulted in a non-nil response`)
	}
	if e == nil {
		t.Fatal(`Nil response and non-nil error resulted in a nil error`)
	}

	// Test that a successful response creates no error
	if _, e = WrapHttpResponseError(&http.Response{StatusCode: http.StatusOK}, nil); e != nil {
		t.Fatal(`Successful response and nil error input resulted in a non-nil err`)
	}

	// Test that a 400 level "Client Error" returns a NonRetriable ClientError and the original response
	ir = &http.Response{StatusCode: http.StatusBadRequest}
	rr, e = WrapHttpResponseError(ir, nil)
	if e == nil {
		t.Fatal(`400 response and non-nil error input resulted in a nil err`)
	}
	if e.IsRetriable() {
		t.Fatal(`400 response and non-nil error input resulted in a Retriable err`)
	}
	if rr == nil || rr != ir {
		t.Fatal(`400 response and non-nil error input failed to return the original response`)
	}

	// Test that a 500 level "Service Fatal" returns a Retriable ClientError and the original response
	ir = &http.Response{StatusCode: http.StatusInternalServerError}
	rr, e = WrapHttpResponseError(ir, nil)
	if e == nil {
		t.Fatal(`500 response and non-nil error input resulted in a nil err`)
	}
	if !e.IsRetriable() {
		t.Fatal(`500 response and non-nil error input resulted in a NonRetriable err`)
	}
	if rr == nil || rr != ir {
		t.Fatal(`500 response and non-nil error input failed to return the original response`)
	}
}
