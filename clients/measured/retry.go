package measured

import (
	//"context"
	"time"
	"github.com/buildertools/svctools-go/clients"
	"github.com/rcrowley/go-metrics"
)

type Collectors struct {
	AttemptCounter metrics.Counter
	ErrorCounter   metrics.Counter
	FatalCounter   metrics.Counter
	TotalTime      metrics.Timer
	AttemptTime    metrics.Timer
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
			TTL: timeout,
			Initial: base,
			MaxJitter: maxJitter,
			Bof: clients.ConstantBackoff,
			Jf: clients.Jitter,
		}, c)
}
func RetryLinear(f clients.RetryFunc,
		timeout time.Duration, 
		base time.Duration, 
		maxJitter time.Duration, 
		c Collectors) (interface{}, error) {
	return Retry(f, &clients.JitteredBackoff{
			TTL: timeout,
			Initial: base,
			MaxJitter: maxJitter,
			Bof: clients.LinearBackoff,
			Jf: clients.Jitter,
		}, c)
}
func RetryExponential(f clients.RetryFunc,
		timeout time.Duration, 
		base time.Duration, 
		maxJitter time.Duration, 
		c Collectors) (interface{}, error) {
	return Retry(f, &clients.JitteredBackoff{
			TTL: timeout,
			Initial: base,
			MaxJitter: maxJitter,
			Bof: clients.ExponentialBackoff,
			Jf: clients.Jitter,
		}, c)
}
