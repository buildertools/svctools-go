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
	"testing"
	"time"
)

func TestNoBackoff(t *testing.T) {
	if NoBackoff(0, time.Duration(0)) != time.Duration(0) {
		t.Fatal(`{0,0} returned non 0`)
	}
	if NoBackoff(0, time.Duration(1)) != time.Duration(0) {
		t.Fatal(`{0,1} returned non 0`)
	}
	if NoBackoff(1, time.Duration(0)) != time.Duration(0) {
		t.Fatal(`{1,0} returned non 0`)
	}
	if NoBackoff(1, time.Duration(1)) != time.Duration(0) {
		t.Fatal(`{1,1} returned non 0`)
	}
}

func TestNoJitter(t *testing.T) {
	if NoJitter(time.Duration(0)) != time.Duration(0) {
		t.Fatal(`{0} returned non 0`)
	}
	if NoJitter(time.Duration(1)) != time.Duration(0) {
		t.Fatal(`{1} returned non 0`)
	}
}

func TestConstantBackoff(t *testing.T) {
	if ConstantBackoff(0, time.Duration(1)) != time.Duration(1) {
		t.Fatalf(`Failed to echo %v with round: 0`, time.Duration(1))
	}
	if ConstantBackoff(1, time.Duration(1)) != time.Duration(1) {
		t.Fatalf(`Failed to echo %v with round: 1`, time.Duration(1))
	}
}

func TestLinearBackoff(t *testing.T) {
	if r := LinearBackoff(0, time.Duration(-1)); r != time.Duration(0) {
		t.Fatalf(`Returned %v instead of %v with round: 0 and initial: %v`, r, time.Duration(0), time.Duration(-1))
	}
	if r := LinearBackoff(0, time.Duration(0)); r != time.Duration(0) {
		t.Fatalf(`Returned %v instead of %v with round: 0 and initial: %v`, r, time.Duration(0), time.Duration(0))
	}
	if r := LinearBackoff(0, time.Duration(1)); r != time.Duration(1) {
		t.Fatalf(`Returned %v instead of %v with round: 0 and initial: %v`, r, time.Duration(1), time.Duration(1))
	}
	if r := LinearBackoff(1, time.Duration(1)); r != time.Duration(2) {
		t.Fatalf(`Returned %v instead of %v with round: 1 and initial: %v`, r, time.Duration(2), time.Duration(1))
	}
	if r := LinearBackoff(2, time.Duration(1)); r != time.Duration(3) {
		t.Fatalf(`Returned %v instead of %v with round: 2 and initial: %v`, r, time.Duration(3), time.Duration(1))
	}
}

func TestExponentialBackoff(t *testing.T) {
	if r := ExponentialBackoff(0, time.Duration(-1)); r != time.Duration(0) {
		t.Fatalf(`Returned %v instead of %v with round: 0 and initial: %v`, r, time.Duration(0), time.Duration(-1))
	}
	if r := ExponentialBackoff(0, time.Duration(1)); r != time.Duration(1) {
		t.Fatalf(`Returned %v instead of %v with round: 0 and initial: %v`, r, time.Duration(1), time.Duration(1))
	}
	if r := ExponentialBackoff(1, time.Duration(1)); r != time.Duration(2) {
		t.Fatalf(`Returned %v instead of %v with round: 1 and initial: %v`, r, time.Duration(2), time.Duration(1))
	}
	if r := ExponentialBackoff(2, time.Duration(1)); r != time.Duration(4) {
		t.Fatalf(`Returned %v instead of %v with round: 2 and initial: %v`, r, time.Duration(4), time.Duration(1))
	}
	if r := ExponentialBackoff(3, time.Duration(1)); r != time.Duration(8) {
		t.Fatalf(`Returned %v instead of %v with round: 3 and initial: %v`, r, time.Duration(8), time.Duration(1))
	}
	if r := ExponentialBackoff(1, time.Duration(1<<62)); r != time.Duration(1<<62) {
		t.Fatalf(`Returned %v instead of %v with round: 1 and initial: %v`, r, time.Duration(1<<62), time.Duration(1<<62))
	}
}

func TestJitter(t *testing.T) {}
