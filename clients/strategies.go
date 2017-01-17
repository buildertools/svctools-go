package clients

import (
	"math/rand"
	"time"
)

type BackoffFunc func(r uint, t time.Duration) time.Duration
type JitterFunc func(t time.Duration) time.Duration

func NoBackoff(round uint, zero time.Duration) time.Duration {
	return time.Duration(0)
}
func NoJitter(none time.Duration) time.Duration {
	return time.Duration(0)
}

func ConstantBackoff(round uint, flat time.Duration) time.Duration {
	return flat
}

func LinearBackoff(round uint, initial time.Duration) time.Duration {
	if initial <= 0 {
		return time.Duration(0)
	}
	return time.Duration(round + 1) * initial
}

func ExponentialBackoff(round uint, initial time.Duration) time.Duration {
	if initial <= 0 {
		return time.Duration(0)
	}
	r := initial << round
	if r < 0 {
		return initial
	}
	return time.Duration(r)
}

func Jitter(max time.Duration) time.Duration {
	if max <= 0 {
		return time.Duration(0)
	}
	return time.Duration(rand.Intn(int(max / time.Millisecond))) * time.Millisecond
}
