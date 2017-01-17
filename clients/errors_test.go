package clients 

import (
	"errors"
//	"net/http"
	"testing"
)

func TestRetriableError(t *testing.T) {
	var e ClientError
	e = RetriableError{ E: nil }
	if !e.IsRetriable() {
		t.Fatal(`RetriableError was notretriable.`)
	}
	if e.Error() != nil {
		t.Fatal(`RetriableError with nil error returned non-nil error.`)
	}
	e = RetriableError{ E: errors.New(`non-nil error`) }
	if !e.IsRetriable() {
		t.Fatal(`RetriableError was notretriable.`)
	}
	if e.Error() == nil {
		t.Fatal(`RetriableError with non-nil error returned nil error.`)
	}
}

func TestNonRetriableError(t *testing.T) {
	var e ClientError
	e = NonRetriableError{ E: nil }
	if e.IsRetriable() {
		t.Fatal(`NonRetriableError was retriable.`)
	}
	if e.Error() != nil {
		t.Fatal(`NonRetriableError with nil error returned non-nil error.`)
	}
	e = NonRetriableError{ E: errors.New(`non-nil error`) }
	if e.IsRetriable() {
		t.Fatal(`NonRetriableError was retriable.`)
	}
	if e.Error() == nil {
		t.Fatal(`NonRetriableError with non-nil error returned nil error.`)
	}
}

