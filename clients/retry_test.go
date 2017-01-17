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
	"testing"
	"time"
)

func TestRetry(t *testing.T) {

	backoffCount := 0
	jitterCount := 0
	jitterSpy := func(i time.Duration) time.Duration {
		jitterCount++
		return time.Duration(0)
	}
	backoffSpy := func(r uint, i time.Duration) time.Duration {
		backoffCount++
		return time.Duration(0)
	}

	fec := 0
	wrappedRe := errors.New("nonretriable")
	f := func() (interface{}, ClientError) {
		switch fec {
		case 0:
			fallthrough
		case 1:
			fec++
			return nil, RetriableError{E: nil}
		case 2:
			fec++
			return nil, NonRetriableError{E: wrappedRe}
		default:
			panic("too many executions")
		}
	}

	pw := &JitteredBackoff{
		TTL:     time.Duration(10) * time.Millisecond,
		Initial: time.Duration(0),
		Jf:      jitterSpy,
		Bof:     backoffSpy,
	}

	r, e := Retry(f, pw)
	if r != nil {
		t.Fatal(`Inexplicably returned a non-nil result`)
	}
	if e == nil {
		t.Fatal(`No error was returned`)
	}
	if e != wrappedRe {
		t.Fatal(`Failed to return the original wrapped error`)
	}
}

// The difficulty in exhaustive testing of WaitOrDie is that functionality depends on timing.
// However, it is simple to test that the injected backoff and jitter functions are used
// and that illegal input is handled correctly.
func TestJitteredBackoffWaitOrDie(t *testing.T) {
	var pw PerishableWaiter
	// TTL, Initial, and MaxJitter are all zero values
	pw = &JitteredBackoff{Jf: NoJitter}
	rc := 0
	recovery := func() {
		if r := recover(); r != nil {
			rc++
		}
	}

	nt := func() {
		defer recovery()
		pw.WaitOrDie(nil)
	}
	nt()
	if rc != 1 {
		t.Fatal(`No panic detected with nil Bof`)
	}

	pw = &JitteredBackoff{Bof: NoBackoff}
	nt = func() {
		defer recovery()
		pw.WaitOrDie(nil)
	}
	nt()
	if rc != 2 {
		t.Fatal(`No panic detected with nil Jf`)
	}

	backoffCount := 0
	jitterCount := 0
	jitterSpy := func(i time.Duration) time.Duration {
		jitterCount++
		return time.Duration(0)
	}
	backoffSpy := func(r uint, i time.Duration) time.Duration {
		backoffCount++
		return time.Duration(0)
	}

	pw = &JitteredBackoff{
		TTL:     time.Duration(0),
		Initial: time.Duration(1) * time.Millisecond,
		Jf:      jitterSpy,
		Bof:     backoffSpy,
	}
	pw.Start()
	ie := errors.New(`e1`)
	e := pw.WaitOrDie(ie)
	if e == nil {
		t.Fatal(`Returned error was nil`)
	}
	if e != ie {
		t.Fatal(`Returned error was not the same object as the input`)
	}
	if backoffCount != 1 || jitterCount != 1 {
		t.Fatal(`backoff or jitter were executed while pw was dead.`)
	}

	backoffCount = 0
	jitterCount = 0

	pw = &JitteredBackoff{
		TTL:     time.Duration(10) * time.Millisecond,
		Initial: time.Duration(0),
		Jf:      jitterSpy,
		Bof:     backoffSpy,
	}
	pw.Start()
	pw.WaitOrDie(nil)
	if backoffCount == 0 {
		t.Fatal(`Backoff not invoked on live`)
	}
	if jitterCount == 0 {
		t.Fatal(`Jitter not invoked on live`)
	}

}

// Validates the mutation caused by (PerishableWaiter) Start()
// Validates the IsDying accessor
// Validates JitteredBackoff implementation of PerishableWaiter
func TestJitteredBackoffStartIsDying(t *testing.T) {
	var pw PerishableWaiter
	pw = &JitteredBackoff{TTL: time.Duration(1)}
	if pw.IsDying() {
		t.Fatal(`Shouldn't be dying yet`)
	}
	pw.Start()
	if !pw.IsDying() {
		t.Fatal(`Should be dying by now`)
	}
}
