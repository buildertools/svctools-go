package clients 

import (
	"testing"
	"time"
)

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
	if r := ExponentialBackoff(1, time.Duration(1 << 62)); r != time.Duration(1 << 62) {
		t.Fatalf(`Returned %v instead of %v with round: 1 and initial: %v`, r, time.Duration(1<<62), time.Duration(1<<62))
	}
}

func TestJitter(t *testing.T) {}

