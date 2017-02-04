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

package measured

import (
	//"context"
	"github.com/buildertools/svctools-go/clients"
	// "github.com/rcrowley/go-metrics"
	"time"
)

type Collectors struct {
	AttemptCounter Counter
	ErrorCounter   Counter
	FatalCounter   Counter
	TotalTime      Timer
	AttemptTime    Timer
}

type Timer interface {
	Update(time.Duration)
}

type Counter interface {
	Inc(int64)
}

func Retry(f clients.RetryFunc, pw clients.PerishableWaiter, c Collectors) (interface{}, error) {
	t0 := time.Now()
	defer func() {
		c.TotalTime.Update(time.Since(t0))
	}()

	pw.Start()
	for {
		c.AttemptCounter.Inc(1)
		tn := time.Now()
		result, err := f()
		c.AttemptTime.Update(time.Since(tn))
		if err == nil {
			return result, nil
		} else if !err.IsRetriable() {
			c.ErrorCounter.Inc(1)
			return nil, err.Error()
		}
		c.FatalCounter.Inc(1)

		if e := pw.WaitOrDie(err.Error()); e != nil {
			return nil, err.Error()
		}
	}
}

func RetryPeriodic(f clients.RetryFunc,
	timeout time.Duration,
	base time.Duration,
	maxJitter time.Duration,
	c Collectors) (interface{}, error) {
	return Retry(f, &clients.JitteredBackoff{
		TTL:       timeout,
		Initial:   base,
		MaxJitter: maxJitter,
		Bof:       clients.ConstantBackoff,
		Jf:        clients.Jitter,
	}, c)
}
func RetryLinear(f clients.RetryFunc,
	timeout time.Duration,
	base time.Duration,
	maxJitter time.Duration,
	c Collectors) (interface{}, error) {
	return Retry(f, &clients.JitteredBackoff{
		TTL:       timeout,
		Initial:   base,
		MaxJitter: maxJitter,
		Bof:       clients.LinearBackoff,
		Jf:        clients.Jitter,
	}, c)
}
func RetryExponential(f clients.RetryFunc,
	timeout time.Duration,
	base time.Duration,
	maxJitter time.Duration,
	c Collectors) (interface{}, error) {
	return Retry(f, &clients.JitteredBackoff{
		TTL:       timeout,
		Initial:   base,
		MaxJitter: maxJitter,
		Bof:       clients.ExponentialBackoff,
		Jf:        clients.Jitter,
	}, c)
}
